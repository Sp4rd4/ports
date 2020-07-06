package grpcserver

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sp4rd4/ports/pkg/proto"
	"google.golang.org/grpc"
)

func (ps *Ports) Serve(lis net.Listener) error {
	ps.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(ps.logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	proto.RegisterPortsServer(ps.grpcServer, ps)
	return ps.grpcServer.Serve(lis)
}

func (ps *Ports) GracefulStop() {
	ps.grpcServer.GracefulStop()
}
