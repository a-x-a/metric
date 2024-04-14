package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"

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

	var dbConn *pgxpool.Pool
	if len(cfg.DatabaseDSN) > 0 {
		dbConn = initDBConn(cfg.DatabaseDSN)
	}

	srv := app.NewServer(dbConn, cfg, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.Run(ctx)

	signal := <-sigint

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, time.Second*5)
	defer cancelShutdown()

	srv.Shutdown(ctxShutdown, signal)
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

func initDBConn(dsn string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("unable to parse DATABASE_URL", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("unable to create connection pool", err)
	}

	return pool
}
