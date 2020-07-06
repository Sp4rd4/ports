package grpcserver

import (
	"context"
	"errors"
	"fmt"

	ptypes "github.com/gogo/protobuf/types"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/proto"
	"github.com/sp4rd4/ports/pkg/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorTag = "grpc"

type PortService interface {
	Save(port *domain.Port) error
	Get(id string) (*domain.Port, error)
}

type Ports struct {
	grpcServer *grpc.Server
	service    PortService
	logger     *zap.Logger
}

func New(srvc PortService, logger *zap.Logger) *Ports {
	return &Ports{service: srvc, logger: logger}
}

func (ps *Ports) Get(ctx context.Context, req *proto.PortRequest) (*proto.Port, error) {
	port, err := ps.service.Get(req.GetId())
	if err != nil {
		ps.logger.Error(fmt.Errorf("[%v] get: %w", errorTag, err).Error())
	}

	return proto.PortDomainToProto(port), convertErrToProto(err)
}

func (ps *Ports) Save(ctx context.Context, req *proto.Port) (*ptypes.Empty, error) {
	err := ps.service.Save(proto.PortProtoToDomain(req))
	if err != nil {
		ps.logger.Error(fmt.Errorf("[%v] save: %w", errorTag, err).Error())
	}

	return &ptypes.Empty{}, convertErrToProto(err)
}

func convertErrToProto(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, domain.ErrNotFound.Error())
	case errors.Is(err, service.ErrPortMissingID):
		return status.Error(codes.InvalidArgument, service.ErrPortMissingID.Error())
	case errors.Is(err, service.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, service.ErrInvalidInput.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
