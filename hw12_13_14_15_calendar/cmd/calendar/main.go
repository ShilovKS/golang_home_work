package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	var storage app.Storage
	if config.Database.Mode == "sql" {
		storage := sqlstorage.New()
		err := storage.Connect(context.Background(), config.Database.DSN(), config.Database.Driver)
		if err != nil {
			logg.Error("failed to connect database: " + err.Error())
			return
		}
		defer storage.Close(context.Background())
	} else {
		storage = memorystorage.New()
	}
	calendar := app.New(logg, storage)
	serverConfig := internalhttp.Config{Host: config.Server.Host, Port: config.Server.Port}
	server := internalhttp.NewServer(logg, calendar, serverConfig)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
