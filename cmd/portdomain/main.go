package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/sp4rd4/ports/pkg/delivery/grpcserver"
	"github.com/sp4rd4/ports/pkg/service"
	"github.com/sp4rd4/ports/pkg/storage/postgres"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"

	// db driver
	_ "github.com/lib/pq"
)

type app struct {
	grpcServer         *grpcserver.Ports
	logger             *zap.Logger
	GRPCPort           string `env:"GRPC_PORT,required"`
	DBHost             string `env:"DATABASE_URL,required"`
	DBMigrationsFolder string `env:"MIGRATIONS_FOLDER" envDefault:"migrations"`
	DBMaxIdleConn      int    `env:"POSTGRES_MAX_IDLE_CONN" envDefault:"100"`
	DBMaxConn          int    `env:"POSTGRES_MAX_CONN" envDefault:"100"`
}

func newApp(logger *zap.Logger) (app, error) {
	appVar := app{}
	if err := env.Parse(&appVar); err != nil {
		return app{}, err
	}

	db, err := sql.Open("postgres", appVar.DBHost)
	if err != nil {
		return app{}, fmt.Errorf("db connect: %w", err)
	}
	dbMigrate, err := sql.Open("postgres", appVar.DBHost)
	if err != nil {
		return app{}, fmt.Errorf("db connect: %w", err)
	}
	db.SetMaxIdleConns(appVar.DBMaxIdleConn)
	db.SetMaxOpenConns(appVar.DBMaxConn)
	storage := postgres.New(db)
	err = storage.Migrate(dbMigrate, appVar.DBMigrationsFolder)
	if err != nil {
		return app{}, fmt.Errorf("postgres migrate: %w", err)
	}

	portService := service.NewPortService(storage)

	server := grpcserver.New(portService, logger)

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
		if err := a.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			a.logger.Fatal(fmt.Errorf("serve: %w", err).Error())
		}
	}()
	a.logger.Info("started grpc portdomain")

	<-ctx.Done()

	a.logger.Info("stopping grpc portdomain")

	a.grpcServer.GracefulStop()

	a.logger.Info("stopped grpc portdomain")
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

	app, err := newApp(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	app.serve(ctx)
}
