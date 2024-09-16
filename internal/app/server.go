package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/handler"
	"github.com/a-x-a/go-metric/internal/logger"
	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	server struct {
		Config     config.ServerConfig
		Storage    storage.Storage
		httpServer *http.Server
		logger     *zap.Logger
	}

	withFileStorage interface {
		Save() error
		Load() error
	}
)

const (
	// logLevel - уровень логирования, по умолчанию info.
	logLevel = "info"
)

var (
	// ErrNotSupportLoadFromFile - хранилище не поддерживает загрузку из файла.
	ErrStorageNotSupportLoadFromFile = errors.New("storage doesn't support loading from file")
)

func NewServer() *server {
	logger := logger.InitLogger(logLevel)
	defer logger.Sync()

	cfg := config.NewServerConfig()

	var dbConn *pgxpool.Pool
	if len(cfg.DatabaseDSN) > 0 {
		poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
		if err != nil {
			logger.Panic("unable to parse DATABASE_URL", zap.Error(err))
		}

		dbConn, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			logger.Panic("unable to create connection pool", zap.Error(err))
		}
	}

	ds := storage.NewDataStorage(dbConn, cfg.FileStoregePath, cfg.StoreInterval, logger)
	ms := metricservice.New(ds, logger)
	rt := handler.NewRouter(ms, logger)
	srv := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: rt,
	}

	return &server{
		Config:     cfg,
		Storage:    ds,
		httpServer: srv,
		logger:     logger,
	}
}

func (s *server) Run(ctx context.Context) {
	if len(s.Config.DatabaseDSN) == 0 && len(s.Config.FileStoregePath) > 0 {
		if s.Config.Restore {
			err := s.loadStorage()
			if err != nil {
				s.logger.Warn("restoring storage", zap.Error(err))
			}
		}

		if s.Config.StoreInterval > 0 {
			go s.saveStorage(ctx)
		}
	}

	s.logger.Info("start http server", zap.String("address", s.Config.ListenAddress))

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Panic("failed to start http server", zap.Error(err))
	}
}

func (s *server) Shutdown(ctx context.Context, signal os.Signal) {
	s.logger.Info("start server shutdown", zap.String("signal", signal.String()))

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Warn("server shutdowning error", zap.Error(err))
	}

	if err := s.Storage.Close(); err != nil {
		s.logger.Error("storage close ", zap.Error(err))
	}

	s.logger.Info("successfully server shutdowning")
}

func (s *server) saveStorage(ctx context.Context) {
	if _, ok := s.Storage.(withFileStorage); !ok {
		s.logger.Debug("storage doesn't support saving to file")
		return
	}

	ticker := time.NewTicker(s.Config.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			func() {
				if err := s.Storage.(withFileStorage).Save(); err != nil {
					s.logger.Error("storage saving error", zap.Error(err))
				}
			}()

		case <-ctx.Done():
			s.logger.Info("shutdown storage saving")
			return
		}
	}
}

func (s *server) loadStorage() error {
	ds, ok := s.Storage.(withFileStorage)
	if !ok {
		return ErrStorageNotSupportLoadFromFile
	}

	if err := ds.Load(); err != nil {
		return err
	}

	return nil
}
