package service_test

import (
	"errors"
	"testing"

	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/service"
	"github.com/stretchr/testify/assert"
)

type MockPortStorage struct {
	err  error
	port *domain.Port
}

func (ms *MockPortStorage) Save(*domain.Port) error {
	return ms.err
}

func (ms *MockPortStorage) Get(id string) (*domain.Port, error) {
	return ms.port, ms.err
}

var testError = errors.New("test")

var examplesSave = []struct {
	name        string
	errStorage  error
	errExpected error
	port        *domain.Port
}{
	{
		name:        "No error",
		errStorage:  nil,
		errExpected: nil,
		port:        &domain.Port{ID: "id", City: "city", Name: "Port"},
	},
	{
		name:        "Test error",
		errStorage:  testError,
		errExpected: testError,
		port:        &domain.Port{ID: "id", City: "city", Name: "Port"},
	},
	{
		name:        "Nil port",
		errStorage:  nil,
		errExpected: service.ErrInvalidInput,
		port:        nil,
	},
	{
		name:        "Incorrect port",
		errStorage:  nil,
		errExpected: service.ErrPortMissingID,
		port:        &domain.Port{City: "city", Name: "Port"},
	},
}

func TestSave(t *testing.T) {
	ms := &MockPortStorage{}
	ps := service.NewPortService(ms)
	for _, ex := range examplesSave {
		ms.err = ex.errStorage
		t.Run(ex.name, func(t *testing.T) {
			err := ps.Save(ex.port)
			assert.True(t, errors.Is(err, ex.errExpected), "Error should be same as expected")
		})
	}
}

var examplesGet = []struct {
	name        string
	errStorage  error
	errExpected error
	id          string
	memory      *domain.Port
	expected    *domain.Port
}{
	{
		name:        "No error",
		errStorage:  nil,
		errExpected: nil,
		id:          "id",
		memory:      &domain.Port{ID: "id", City: "city", Name: "Port"},
		expected:    &domain.Port{ID: "id", City: "city", Name: "Port"},
	},
	{
		name:        "Test error",
		errStorage:  testError,
		errExpected: testError,
		id:          "id",
		memory:      &domain.Port{ID: "id", City: "city", Name: "Port"},
		expected:    nil,
	},
	{
		name:        "No id",
		errStorage:  nil,
		errExpected: service.ErrPortMissingID,
		id:          "",
		memory:      &domain.Port{ID: "id", City: "city", Name: "Port"},
		expected:    nil,
	},
}

func TestGet(t *testing.T) {
	ms := &MockPortStorage{}
	ps := service.NewPortService(ms)
	for _, ex := range examplesGet {
		ms.err = ex.errStorage
		ms.port = ex.memory
		t.Run(ex.name, func(t *testing.T) {
			port, err := ps.Get(ex.id)
			assert.True(t, errors.Is(err, ex.errExpected), "Error should be same as expected")
			assert.Equal(t, ex.expected, port, "Should return port same as expected")
		})
	}
}
