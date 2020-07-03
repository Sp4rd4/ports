package loader

import "github.com/sp4rd4/ports/pkg/domain"

type Ports interface {
	Load() <-chan *domain.Port
}
