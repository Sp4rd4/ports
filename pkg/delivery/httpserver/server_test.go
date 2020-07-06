package httpserver_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/sp4rd4/ports/pkg/delivery/httpserver"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type mockService struct {
	err  error
	port *domain.Port
}

func (ms *mockService) Get(id string) (*domain.Port, error) {
	return ms.port, ms.err
}

var (
	testError = errors.New("test")

	json = jsoniter.ConfigDefault
)

var examplesGet = []struct {
	name       string
	id         string
	status     int
	errService error
	errLogged  error
	port       *domain.Port
}{
	{
		name:       "No error",
		errService: nil,
		errLogged:  nil,
		status:     http.StatusOK,
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
		status:     http.StatusInternalServerError,
		port:       nil,
	},
	{
		name:       "No port",
		errService: domain.ErrNotFound,
		errLogged:  nil,
		status:     http.StatusNotFound,
		id:         "AEAJM",
		port:       nil,
	},
}

func TestGet(t *testing.T) {
	ms := &mockService{}
	core, observed := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	handler := httpserver.New(ms, logger)
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	for _, ex := range examplesGet {
		ms.port = ex.port
		ms.err = ex.errService
		observed.TakeAll()

		t.Run(ex.name, func(t *testing.T) {
			expct := e.GET("/ports/" + ex.id).Expect().Status(ex.status)
			if ex.errService != nil {
				expct.JSON().Object().ValueEqual("message", http.StatusText(ex.status))
			}
			if ex.errLogged != nil {
				assert.Equal(t,
					observed.FilterMessage(fmt.Errorf("[http]: %w", ex.errLogged).Error()).Len(), 1,
					"Should log internal server errors",
				)
			}

			if ex.port != nil {
				expct.JSON().Object().Equal(ex.port)
			}
		})
	}
}
