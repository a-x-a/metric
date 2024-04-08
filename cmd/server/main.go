package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/app"
	"github.com/a-x-a/go-metric/internal/config"
)

const (
	// logLevel - уровень логирования, по умолчанию info.
	logLevel = "info"
)

func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	logger := initLogger(logLevel)
	defer logger.Sync()

	cfg := config.NewServerConfig()
	srv := app.NewServer(cfg, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.Run(ctx)

	signal := <-sigint

	srv.Shutdown(ctx, signal)
}

func initLogger(level string) *zap.Logger {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	return zl
}
