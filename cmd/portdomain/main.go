package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sp4rd4/ports/pkg/delivery/grpc"
	"github.com/sp4rd4/ports/pkg/service"
	"github.com/sp4rd4/ports/pkg/storage/postgres"
	"go.uber.org/zap"
	ngrpc "google.golang.org/grpc"
)

const (
	shutdownTimeout = 4 * time.Second
)

type app struct {
	grpcServer *grpc.PortServer
	logger     *zap.Logger
	GRPCPort   string `env:"GRPC_PORT,required"`
}

func newApp(ctx context.Context, logger *zap.Logger) (app, error) {
	appVar := app{}
	if err := env.Parse(&appVar); err != nil {
		return app{}, err
	}

	storageConfig := postgres.Config{}
	if err := env.Parse(&storageConfig); err != nil {
		return app{}, fmt.Errorf("grpc client config read: %w", err)
	}
	storage, err := postgres.New(storageConfig)
	if err != nil {
		return app{}, fmt.Errorf("grpc client init: %w", err)
	}

	portService := service.NewPortService(storage)

	server := grpc.New(portService, logger)

	appVar.grpcServer = server
	appVar.logger = logger

	return appVar, nil
}

func (a *app) serve(ctx context.Context) {
	lis, err := net.Listen("tcp", ":"+a.GRPCPort)
	if err != nil {
		a.logger.Fatal(fmt.Errorf("tcp listen: %w", err).Error())
	}

	go func() {
		if err := a.grpcServer.Serve(lis); err != nil && !errors.Is(err, ngrpc.ErrServerStopped) {
			a.logger.Fatal(fmt.Errorf("serve: %w", err).Error())
		}
	}()
	a.logger.Info("started grpc portsdomain")

	<-ctx.Done()

	a.logger.Info("stopping grpc portsdomain")

	a.grpcServer.GracefulStop()

	a.logger.Info("stopped grpc portsdomain")
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-signals
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	app, err := newApp(ctx, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	app.serve(ctx)
}
