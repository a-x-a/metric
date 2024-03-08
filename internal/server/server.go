package server

import (
	"fmt"
	"net/http"

	"github.com/a-x-a/go-metric/internal/handler"
	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/a-x-a/go-metric/internal/storage/memory"
)

type Server interface {
	Run() error
}

// type metricService interface {
// 	Save(metric string, metricType string, value string) error
// }

// type server struct {
// 	service metricService
// 	storage storage.Storage
// }

func Run() error {
	stor := memory.NewMemStorage()
	service := metricservice.New(stor)
	updateHandler := handler.NewUpdateHandler(service)

	mux := http.NewServeMux()
	mux.Handle("/update/", updateHandler)

	fmt.Println("listening on 8080")

	err := http.ListenAndServe("localhost:8080", mux)

	if err != nil {
		return err
	}

	return nil
}
