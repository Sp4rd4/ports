package service

import (
	"fmt"

	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/domain/loader"
	"go.uber.org/zap"
)

const errorTagLoader = "load-service"

type LoadService struct {
	loader  loader.Ports
	storage domain.PortRepository
	logger  *zap.Logger
}

func NewLoadService(loader loader.Ports, storage domain.PortRepository, logger *zap.Logger) LoadService {
	return LoadService{storage: storage, loader: loader}
}

func (s LoadService) Load() {
	ports := s.loader.Load()
	for port := range ports {
		err := s.storage.Save(port)
		if err != nil {
			s.logger.Error(fmt.Errorf("[%v] loader: %w", errorTagLoader, err).Error())
		}
	}
}
