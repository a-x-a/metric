package handler

import (
	"net/http"

	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/go-chi/chi"
)

func Router(metricService metricservice.MetricService) http.Handler {
	metricHendlers := newMetricHandlers(metricService)

	r := chi.NewRouter()

	// r.Use(logging.RequestsLogger)
	// r.Get("/", metricHendlers.List)
	// r.Get("/value/{kind}/{name}", metricHendlers.Get)
	r.Post("/update/{kind}/{name}/{value}", metricHendlers.Update)

	return r
}
