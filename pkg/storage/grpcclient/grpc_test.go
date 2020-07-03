package grpcclient_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/proto"
	"github.com/sp4rd4/ports/pkg/storage/grpcclient"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockPortsClient struct {
	err          error
	memory       *domain.Port
	grpcResponse *proto.Port
}

func (c *MockPortsClient) Save(_ context.Context, _ *proto.Port, _ ...grpc.CallOption) (*types.Empty, error) {
	return &types.Empty{}, c.err
}

func (c *MockPortsClient) Get(_ context.Context, in *proto.PortRequest, _ ...grpc.CallOption) (*proto.Port, error) {
	if c.err != nil {
		return &proto.Port{}, c.err
	}

	if c.memory != nil && in.Id == c.memory.ID {
		return c.grpcResponse, nil
	}
	return &proto.Port{}, status.Error(codes.NotFound, domain.ErrNotFound.Error())
}

type GRPCTestSuite struct {
	suite.Suite
	mock    *MockPortsClient
	storage domain.PortRepository
}

func (s *GRPCTestSuite) SetupSuite() {
	s.mock = &MockPortsClient{}
	s.storage = grpcclient.New(s.mock)
}

var testError = errors.New("test")

var examplesSave = []struct {
	name   string
	errSet error
	errGot error
}{
	{
		name:   "No error",
		errSet: nil,
		errGot: nil,
	},
	{
		name:   "Test error",
		errSet: testError,
		errGot: testError,
	},
}

func (s *GRPCTestSuite) TestSave() {
	for _, ex := range examplesSave {
		s.mock.err = ex.errSet
		s.Run(ex.name, func() {
			err := s.storage.Save(nil)
			s.True(errors.Is(err, ex.errGot), "Error should be same as expected")
		})
	}
}

var examplesGet = []struct {
	name         string
	errSet       error
	errGot       error
	id           string
	memory       *domain.Port
	grpcResponse *proto.Port
	expected     *domain.Port
}{
	{
		name:         "No error",
		errSet:       nil,
		errGot:       nil,
		id:           "id",
		memory:       &domain.Port{ID: "id", City: "city", Name: "Port"},
		grpcResponse: &proto.Port{Id: "id", City: "city", Name: "Port"},
		expected:     &domain.Port{ID: "id", City: "city", Name: "Port"},
	},
	{
		name:         "Test error",
		errSet:       testError,
		errGot:       testError,
		memory:       &domain.Port{ID: "id", City: "city", Name: "Port"},
		grpcResponse: nil,
		expected:     nil,
	},
	{
		name:         "Not found",
		errSet:       status.Error(codes.NotFound, domain.ErrNotFound.Error()),
		errGot:       domain.ErrNotFound,
		id:           "notid",
		memory:       &domain.Port{ID: "id", City: "city", Name: "Port"},
		grpcResponse: &proto.Port{},
		expected:     nil,
	},
}

func (s *GRPCTestSuite) TestGet() {
	for _, ex := range examplesGet {
		s.mock.err = ex.errSet
		s.mock.memory = ex.memory
		s.mock.grpcResponse = ex.grpcResponse
		s.Run(ex.name, func() {
			port, err := s.storage.Get(ex.id)
			s.True(errors.Is(err, ex.errGot), "Error should be same as expected")
			s.Equal(ex.expected, port, "Should return port same as expected")
		})
	}
}

func TestGRPCTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}
