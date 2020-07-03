package service

import (
	"errors"
	"fmt"

	"github.com/sp4rd4/ports/pkg/domain"
)

const errorTagPort = "load-service"

var (
	ErrPortMissingID = errors.New("port missing id")
	ErrInvalidInput  = errors.New("invalid input")
)

type PortService struct {
	storage domain.PortRepository
}

func NewPortService(storage domain.PortRepository) PortService {
	return PortService{storage: storage}
}

func (s PortService) Save(port *domain.Port) error {
	if port == nil {
		return fmt.Errorf("[%v] save: %w", errorTagPort, ErrInvalidInput)
	}
	if port.ID == "" {
		return fmt.Errorf("[%v] save: %w", errorTagPort, ErrPortMissingID)
	}
	err := s.storage.Save(port)
	if err != nil {
		return fmt.Errorf("[%v] save: %w", errorTagPort, err)
	}
	return nil
}

func (s PortService) Get(id string) (*domain.Port, error) {
	if id == "" {
		return nil, fmt.Errorf("[%v] get: %w", errorTagPort, ErrPortMissingID)
	}
	port, err := s.storage.Get(id)
	if err != nil {
		return nil, fmt.Errorf("[%v] get: %w", errorTagPort, err)
	}
	return port, nil
}
