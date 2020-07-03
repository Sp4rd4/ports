package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	nhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/machinebox/progress"
	"github.com/sp4rd4/ports/pkg/delivery/http"
	"github.com/sp4rd4/ports/pkg/jsonreader"
	"github.com/sp4rd4/ports/pkg/proto"
	"github.com/sp4rd4/ports/pkg/service"
	"github.com/sp4rd4/ports/pkg/storage/grpcclient"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	loadNotifyInterval = 10 * time.Second
	shutdownTimeout    = 4 * time.Second
)

type app struct {
	server           *http.PortController
	loadService      *service.LoadService
	logger           *zap.Logger
	PortsFilepath    string        `env:"PORTS_FILE,required"`
	HTTPPort         string        `env:"HTTP_PORT,required"`
	HTTPReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
	HTTPWriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	HTTPIdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"120s"`
	PortDomainHost   string        `env:"PORTS_DOMAIN_HOST,required"`
	LoaderBufferSize int           `env:"JSON_BUFFER_SIZE" envDefault:"512"`
	PoolSize         int           `env:"WORKER_POOL_SIZE" envDefault:"100"`
}

func newApp(ctx context.Context, logger *zap.Logger) (app, error) {
	appVar := app{}
	if err := env.Parse(&appVar); err != nil {
		return app{}, err
	}

	conn, err := grpc.Dial(appVar.PortDomainHost, grpc.WithInsecure())
	if err != nil {
		return app{}, fmt.Errorf("portdomain connect: %w", err)
	}
	storage := grpcclient.New(proto.NewPortsClient(conn))

	info, err := os.Stat(appVar.PortsFilepath)
	if err != nil {
		return app{}, fmt.Errorf("get file info: %w", err)
	}
	size := info.Size()
	file, err := os.Open(appVar.PortsFilepath)
	if err != nil {
		return app{}, fmt.Errorf("open file: %w", err)
	}

	portService := service.NewPortService(storage)

	controller := http.New(portService, logger)

	appVar.server = controller
	appVar.logger = logger

	loader := jsonreader.NewLoader(meterReader(file, size), appVar.LoaderBufferSize, ctx.Done())
	loadService, err := service.NewLoadService(loader, storage, logger, appVar.PoolSize)
	if err != nil {
		return app{}, fmt.Errorf("start service: %w", err)
	}

	appVar.loadService = &loadService

	return appVar, nil
}

func (a *app) serve(ctx context.Context) {
	srv := &nhttp.Server{
		ReadTimeout:  a.HTTPReadTimeout,
		WriteTimeout: a.HTTPWriteTimeout,
		IdleTimeout:  a.HTTPIdleTimeout,
		Handler:      a.server,
		Addr:         ":" + a.HTTPPort,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, nhttp.ErrServerClosed) {
			a.logger.Fatal(fmt.Errorf("serve: %w", err).Error())
		}
	}()
	a.logger.Info("started http clientapi")

	<-ctx.Done()

	a.logger.Info("stopping http clientapi")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		a.logger.Fatal(fmt.Errorf("shutdown: %w", err).Error())
	}

	a.logger.Info("stopped http clientapi")
}

func meterReader(r io.Reader, size int64) io.Reader {
	meteredReader := progress.NewReader(r)

	// Start a goroutine printing progress
	go func() {
		ctx := context.Background()
		progressChan := progress.NewTicker(ctx, meteredReader, size, loadNotifyInterval)
		for p := range progressChan {
			fmt.Printf("%.2f%% of json processed...\n", p.Percent())
		}
		fmt.Println("json is  processed.")
	}()

	return meteredReader
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

	go app.loadService.Load()

	app.serve(ctx)
}
