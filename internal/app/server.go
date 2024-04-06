package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

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
	}

	withFileStorage interface {
		Save() error
		Load() error
	}
)

func NewServer() *server {
	cfg := config.NewServerConfig()
	ds := storage.NewDataStorage(cfg.FileStoregePath, cfg.StoreInterval)
	ms := metricservice.New(ds)
	rt := handler.Router(ms)
	srv := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: rt,
	}

	return &server{
		Config:     cfg,
		Storage:    ds,
		httpServer: srv,
	}
}

func (s *server) Run(ctx context.Context) {
	if err := logger.Initialize(s.Config.LogLevel); err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	defer logger.Log.Sync()

	if s.Config.Restore {
		s.loadStorage()
	}

	if len(s.Config.FileStoregePath) > 0 && s.Config.StoreInterval > 0 {
		go s.saveStorage(ctx)
	}

	logger.Log.Info("start http server", zap.String("address", s.Config.ListenAddress))

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log.Fatal("failed to start http server", zap.Error(err))
	}
}

func (s *server) Shutdown(ctx context.Context, signal os.Signal) {
	logger.Log.Info("start server shutdown", zap.String("signal", signal.String()))

	ctxShutdown, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := s.httpServer.Shutdown(ctxShutdown); err != nil {
		logger.Log.Error("server shutdowning error", zap.Error(err))
	}

	if ds, ok := s.Storage.(withFileStorage); ok {
		if err := ds.Save(); err != nil {
			logger.Log.Error("storage saving error", zap.Error(err))
		}
	}

	logger.Log.Info("successfully server shutdowning")
}

func (s *server) saveStorage(ctx context.Context) {
	if _, ok := s.Storage.(withFileStorage); !ok {
		logger.Log.Debug("storage doesn't support saving to file")
		return
	}

	ticker := time.NewTicker(s.Config.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			func() {
				if err := s.Storage.(withFileStorage).Save(); err != nil {
					logger.Log.Error("storage saving error", zap.Error(err))
				}
			}()

		case <-ctx.Done():
			logger.Log.Info("shutdown storage saving")
			return
		}
	}
}

func (s *server) loadStorage() {
	ds, ok := s.Storage.(withFileStorage)
	if !ok {
		logger.Log.Debug("storage doesn't support loading from file")
		return
	}

	if err := ds.Load(); err != nil {
		logger.Log.Panic("failed to load storage", zap.Error(err))
	}
}
