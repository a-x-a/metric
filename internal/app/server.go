package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/handler"
	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	Server interface {
		Run() error
	}

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
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("listening on", s.Config.ListenAddress)
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start http server: %v", err))
		}
	}()

	wg.Wait()
}
