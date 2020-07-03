package grpcclient

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorTag = "grpc-store"

type storage struct {
	client proto.PortsClient
}

type Config struct {
	PortDomainHost string `env:"PORTS_DOMAIN_HOST,required"`
}

func New(conf Config) (domain.PortRepository, error) {
	conn, err := grpc.Dial(conf.PortDomainHost, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("[%v] portdomain connect: %w", errorTag, err)
	}

	return &storage{client: proto.NewPortsClient(conn)}, nil
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

	return proto.PortProtoToDomain(port), fmt.Errorf("[%v] save: %w", errorTag, err)
}
