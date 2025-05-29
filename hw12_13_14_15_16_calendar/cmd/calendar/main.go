package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Calendar/hw12_13_14_15_calendar/internal/app"
	"github.com/Calendar/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Calendar/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Calendar/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Calendar/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	fmt.Printf("Database config: %+v\n", config)
	var storage app.Storage
	switch config.Storage {
	case "database":
		storage = sqlstorage.New(config.Database)
	case "in-memory":
		storage = memorystorage.New()
	default:
		storage = memorystorage.New()
	}
	Logg := logger.New(config.Logger.Level)
	calendar := app.New(Logg, storage)
	server := internalhttp.NewServer(Logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			Logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	Logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		Logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
