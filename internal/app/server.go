package app

import (
	"fmt"
	"net/http"

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
		Config  config.ServerConfig
		Storage storage.Storage
	}
)

func NewServer() *server {
	return &server{
		Config:  config.NewServerConfig(),
		Storage: storage.NewMemStorage(),
	}
}

func (s *server) Run() error {
	service := metricservice.New(s.Storage)
	r := handler.Router(service)

	fmt.Println("listening on", s.Config.ListenAddress)

	return http.ListenAndServe(s.Config.ListenAddress, r)
}
