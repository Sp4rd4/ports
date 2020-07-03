package grpcclient

import (
	"context"
	"fmt"

	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorTag = "grpc-store"

type storage struct {
	client proto.PortsClient
}

func New(client proto.PortsClient) domain.PortRepository {
	return &storage{client: client}
}

func (s storage) Save(port *domain.Port) error {
	_, err := s.client.Save(context.Background(), proto.PortDomainToProto(port))
	if err != nil {
		return fmt.Errorf("[%v] save: %w", errorTag, err)
	}
	return nil
}

func (s storage) Get(id string) (*domain.Port, error) {
	port, err := s.client.Get(context.Background(), &proto.PortRequest{Id: id})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return nil, fmt.Errorf("[%v] save: %w", errorTag, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("[%v] save: %w", errorTag, err)
	}

	return proto.PortProtoToDomain(port), nil
}
