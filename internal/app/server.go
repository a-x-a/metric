package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"

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
)

func NewServer() *server {
	cfg := config.NewServerConfig()
	ds := storage.NewMemStorage()
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

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		logger.Log.Info("start http server", zap.String("address", s.Config.ListenAddress))
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Panic("failed to start http server", zap.String("err", err.Error()))
		}
	}()

	wg.Wait()
}
