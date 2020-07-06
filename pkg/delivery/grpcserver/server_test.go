package grpcserver_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sp4rd4/ports/pkg/delivery/grpcserver"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/proto"
	"github.com/sp4rd4/ports/pkg/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockService struct {
	err  error
	port *domain.Port
}

func (ms *mockService) Get(id string) (*domain.Port, error) {
	return ms.port, ms.err
}
func (ms *mockService) Save(p *domain.Port) error {
	ms.port = p
	return ms.err
}

var (
	testError = errors.New("test")
)

type GRPCTestSuite struct {
	suite.Suite
	mock     *mockService
	server   *grpcserver.Ports
	logger   *zap.Logger
	observed *observer.ObservedLogs
}

func (s *GRPCTestSuite) SetupSuite() {
	s.mock = &mockService{}
	core, observed := observer.New(zapcore.DebugLevel)
	s.logger = zap.New(core)
	s.observed = observed
	s.server = grpcserver.New(s.mock, s.logger)
}

var examplesGet = []struct {
	name       string
	id         string
	status     codes.Code
	errService error
	errLogged  error
	port       *domain.Port
}{
	{
		name:       "No error",
		errService: nil,
		errLogged:  nil,
		id:         "AEAJM",
		port: &domain.Port{
			ID:      "AEAJM",
			Name:    "Ajman",
			City:    "Ajman",
			Country: "United Arab Emirates",
			Alias:   domain.StringArray{},
			Regions: domain.StringArray{},
			Coordinates: domain.Location{
				Latitude:  55.5136433,
				Longitude: 25.4052165,
			},
			Province: "Ajman",
			Timezone: "Asia/Dubai",
			Unlocs:   domain.StringArray{"AEAJM"},
			Code:     "52000",
		},
	},
	{
		name:       "Test error",
		errService: testError,
		errLogged:  testError,
		id:         "AEAJM",
		status:     codes.Internal,
		port:       nil,
	},
	{
		name:       "Not found",
		errService: domain.ErrNotFound,
		errLogged:  domain.ErrNotFound,
		status:     codes.NotFound,
		id:         "AEAJM",
		port:       nil,
	},
	{
		name:       "Missing ID",
		errService: service.ErrPortMissingID,
		errLogged:  service.ErrPortMissingID,
		status:     codes.InvalidArgument,
		id:         "",
		port:       nil,
	},
}

func (s *GRPCTestSuite) TestGet() {
	for _, ex := range examplesGet {
		s.mock.port = ex.port
		s.mock.err = ex.errService

		s.observed.TakeAll()

		s.Run(ex.name, func() {
			port, err := s.server.Get(context.TODO(), &proto.PortRequest{Id: ex.id})
			s.Equal(proto.PortDomainToProto(ex.port), port, "Should return expected port")

			if err != nil {
				st := status.Convert(err)
				s.Equal(ex.status, st.Code(), "Should return expected error code")
				s.Equal(ex.errLogged.Error(), st.Message(), "Should return expected error message")
				s.Equal(
					1, s.observed.FilterMessage(fmt.Errorf("[grpc] get: %w", ex.errLogged).Error()).Len(),
					"Should contain appropriate log message",
				)
			}
		})
	}
}

var examplesSave = []struct {
	name       string
	id         string
	status     codes.Code
	errService error
	errLogged  error
	port       *domain.Port
}{
	{
		name:       "No error",
		errService: nil,
		errLogged:  nil,
		id:         "AEAJM",
		port: &domain.Port{
			ID:      "AEAJM",
			Name:    "Ajman",
			City:    "Ajman",
			Country: "United Arab Emirates",
			Alias:   domain.StringArray{},
			Regions: domain.StringArray{},
			Coordinates: domain.Location{
				Latitude:  55.5136433,
				Longitude: 25.4052165,
			},
			Province: "Ajman",
			Timezone: "Asia/Dubai",
			Unlocs:   domain.StringArray{"AEAJM"},
			Code:     "52000",
		},
	},
	{
		name:       "Test error",
		errService: testError,
		errLogged:  testError,
		id:         "AEAJM",
		status:     codes.Internal,
		port:       &domain.Port{},
	},
	{
		name:       "Missing ID",
		errService: service.ErrPortMissingID,
		errLogged:  service.ErrPortMissingID,
		status:     codes.InvalidArgument,
		id:         "",
		port:       &domain.Port{},
	},
	{
		name:       "Invalid input",
		errService: service.ErrInvalidInput,
		errLogged:  service.ErrInvalidInput,
		status:     codes.InvalidArgument,
		id:         "",
		port:       &domain.Port{},
	},
}

func (s *GRPCTestSuite) TestSave() {
	for _, ex := range examplesSave {
		s.mock.port = nil
		s.mock.err = ex.errService
		s.observed.TakeAll()
		s.Run(ex.name, func() {
			_, err := s.server.Save(context.TODO(), proto.PortDomainToProto(ex.port))
			s.Equal(ex.port, s.mock.port, "Should save expected port")

			if err != nil {
				st := status.Convert(err)
				s.Equal(ex.status, st.Code(), "Should return expected error code")
				s.Equal(ex.errLogged.Error(), st.Message(), "Should return expected error message")
				s.Equal(
					1, s.observed.FilterMessage(fmt.Errorf("[grpc] save: %w", ex.errLogged).Error()).Len(),
					"Should contain appropriate log message",
				)
			}
		})
	}
}

func TestGRPCTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}
