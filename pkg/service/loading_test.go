package service_test

import (
	"strconv"
	"testing"

	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type MockPortSliceStorage struct {
	err   error
	ports []*domain.Port
}

func (ms *MockPortSliceStorage) Save(p *domain.Port) error {
	if ms.err != nil {
		return ms.err
	}
	ms.ports = append(ms.ports, p)
	return nil
}

func (ms *MockPortSliceStorage) Get(id string) (*domain.Port, error) {
	return nil, nil
}

type loaderSlice struct {
	ports []*domain.Port
}

func (l *loaderSlice) Load() <-chan *domain.Port {
	res := make(chan *domain.Port)
	go func() {
		for _, p := range l.ports {
			res <- p
		}
		close(res)
	}()
	return res
}

var examplesLoad = []struct {
	name       string
	errStorage error
	ports      []*domain.Port
}{
	{
		name:       "No error",
		errStorage: nil,
		ports: []*domain.Port{
			{
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
			{
				ID:      "ZAPLZ",
				Name:    "Port Elizabeth",
				City:    "Port Elizabeth",
				Country: "South Africa",
				Alias:   domain.StringArray{},
				Regions: domain.StringArray{},
				Coordinates: domain.Location{
					Latitude:  25.5207358,
					Longitude: -33.7139247,
				},
				Province: "Eastern Cape",
				Timezone: "Africa/Johannesburg",
				Unlocs:   domain.StringArray{"ZAPLZ"},
				Code:     "79145",
			},
		},
	},
	{
		name:       "Test error",
		errStorage: errFoo,
		ports:      nil,
	},
	{
		name:       "No ports",
		errStorage: nil,
		ports:      nil,
	},
	{
		name:       "Empty ports",
		errStorage: nil,
		ports:      []*domain.Port{},
	},
}

func TestLoad(t *testing.T) {
	ms := &MockPortSliceStorage{}
	core, observed := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	ldr := &loaderSlice{}
	ps, err := service.NewLoadService(ldr, ms, logger, 100)
	if err != nil {
		t.Fatalf("Could not create service: %s", err)
	}

	for _, ex := range examplesLoad {
		ldr.ports = ex.ports
		ms.err = ex.errStorage
		ms.ports = nil
		observed.TakeAll()

		t.Run(ex.name, func(t *testing.T) {
			ps.Load()
			if ex.errStorage == nil {
				assert.Zero(t, observed.Len(), "Should be zero errors logged")
			} else {
				assert.Equalf(t, observed.Len(), len(ex.ports), "Should be %v errors logged", strconv.Itoa(len(ex.ports)))
			}
			assert.ElementsMatch(t, ms.ports, ex.ports, "Ports should be same as expected")
		})
	}
}
