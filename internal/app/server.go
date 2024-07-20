// Package app инициализация и запуск приложения.
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
	"github.com/a-x-a/go-metric/internal/security"
	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	Server struct {
		config     config.ServerConfig
		storage    storage.Storage
		httpServer *http.Server
		logger     *zap.Logger
		key        security.PrivateKey
	}

	WithFileStorage interface {
		Save() error
		Load() error
	}
)

var (
	// ErrNotSupportLoadFromFile хранилище не поддерживает загрузку из файла.
	ErrStorageNotSupportLoadFromFile = errors.New("storage doesn't support loading from file")
)

func NewServer(logLevel string) *Server {
	var err error
	log := logger.InitLogger(logLevel)
	defer log.Sync()

	cfg := config.NewServerConfig()
	err = cfg.Parse()
	if err != nil {
		log.Panic("server failed to parse config", zap.Error(err))
	}

	var privateKey security.PrivateKey
	if len(cfg.CryptoKey) != 0 {
		privateKey, err = security.NewPrivateKey(cfg.CryptoKey)
		if err != nil {
			log.Panic("server failed to get private key", zap.Error(err))
		}
	}

	var dbConn *pgxpool.Pool
	if len(cfg.DatabaseDSN) > 0 {
		poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
		if err != nil {
			log.Panic("unable to parse DATABASE_URL", zap.Error(err))
		}

		dbConn, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Panic("unable to create connection pool", zap.Error(err))
		}

		if err := migrationRun(cfg.DatabaseDSN, log); err != nil {
			log.Panic("unable to init DB", zap.Error(err))
		}
	}

	ds := storage.NewDataStorage(dbConn, cfg.FileStoregePath, cfg.StoreInterval, log)
	ms := metricservice.New(ds, log)
	rt := handler.NewRouter(ms, log, cfg.Key, privateKey)
	srv := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: rt,
	}

	return &Server{
		config:     cfg,
		storage:    ds,
		httpServer: srv,
		logger:     log,
		key:        privateKey,
	}
}

func (s *Server) Run(ctx context.Context) {
	if len(s.config.DatabaseDSN) == 0 && len(s.config.FileStoregePath) > 0 {
		if s.config.Restore {
			err := s.loadStorage()
			if err != nil {
				s.logger.Warn("restoring storage", zap.Error(err))
			}
		}

		if s.config.StoreInterval > 0 {
			go s.saveStorage(ctx)
		}
	}

	s.logger.Info("start http server", zap.String("address", s.config.ListenAddress))

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Panic("failed to start http server", zap.Error(err))
	}
}

func (s *Server) Shutdown(ctx context.Context, signal os.Signal) {
	s.logger.Info("start server shutdown", zap.String("signal", signal.String()))

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Warn("server shutdowning error", zap.Error(err))
	}

	if err := s.storage.Close(); err != nil {
		s.logger.Error("storage close ", zap.Error(err))
	}

	s.logger.Info("successfully server shutdowning")
}

func (s *Server) saveStorage(ctx context.Context) {
	if _, ok := s.storage.(WithFileStorage); !ok {
		s.logger.Debug("storage doesn't support saving to file")
		return
	}

	ticker := time.NewTicker(s.config.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			func() {
				if err := s.storage.(WithFileStorage).Save(); err != nil {
					s.logger.Error("storage saving error", zap.Error(err))
				}
			}()

		case <-ctx.Done():
			s.logger.Info("shutdown storage saving")
			return
		}
	}
}

func (s *Server) loadStorage() error {
	ds, ok := s.storage.(WithFileStorage)
	if !ok {
		return ErrStorageNotSupportLoadFromFile
	}

	if err := ds.Load(); err != nil {
		return err
	}

	return nil
}
