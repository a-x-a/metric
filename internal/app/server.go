package app

import (
	"fmt"
	"net/http"

	"github.com/a-x-a/go-metric/internal/handler"
	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	Server interface {
		Run() error
	}
	serverConfig struct {
		// ListenAddress - адрес сервера сбора метрик
		ListenAddress string
	}
	server struct {
		Config serverConfig
		// Service MetricService
		// Storage storage.Storage
	}
)

// type metricService interface {
// 	Save(metric string, metricType string, value string) error
// }

// type server struct {
// 	service metricService
// 	storage storage.Storage
// }

func newServerConfig() serverConfig {
	return serverConfig{
		ListenAddress: "localhost:8080",
	}
}

func NewServer() *server {
	return &server{Config: newServerConfig()}
}

func (s *server) Run() error {
	stor := storage.NewMemStorage()
	service := metricservice.New(stor)
	updateHandler := handler.NewUpdateHandler(service)
	mux := http.NewServeMux()
	mux.Handle("/update/", updateHandler)

	fmt.Println("listening on", s.Config.ListenAddress)

	err := http.ListenAndServe(s.Config.ListenAddress, mux)

	if err != nil {
		return err
	}

	return nil
}
