package service

import (
	"fmt"

	"github.com/panjf2000/ants/v2"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/domain/loader"
	"go.uber.org/zap"
)

const errorTagLoader = "load-service"

type LoadService struct {
	loader  loader.Ports
	storage domain.PortRepository
	logger  *zap.Logger
	pool    *ants.Pool
}

func NewLoadService(
	ldr loader.Ports, storage domain.PortRepository, logger *zap.Logger, workers int,
) (LoadService, error) {
	pool, err := ants.NewPool(workers)
	if err != nil {
		return LoadService{}, fmt.Errorf("[%v] pool init: %w", errorTagLoader, err)
	}
	return LoadService{storage: storage, loader: ldr, logger: logger, pool: pool}, nil
}

func (s LoadService) Load() {
	ports := s.loader.Load()
	for p := range ports {
		lp := p
		err := ants.Submit(func() {
			err := s.storage.Save(lp)
			if err != nil {
				s.logger.Error(fmt.Errorf("[%v] save: %w", errorTagLoader, err).Error())
			}
		})
		if err != nil {
			s.logger.Error(fmt.Errorf("[%v] pool submit: %w", errorTagLoader, err).Error())
		}
	}
}
