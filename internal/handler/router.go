package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/a-x-a/go-metric/internal/logger"
)

func Router(s metricService) http.Handler {
	metricHendlers := newMetricHandlers(s)
	r := chi.NewRouter()
	r.Use(logger.WithLogger)

	r.Get("/", metricHendlers.List)
	r.Get("/value/{kind}/{name}", metricHendlers.Get)
	r.Post("/update/{kind}/{name}/{value}", metricHendlers.Update)

	r.Post("/value/", metricHendlers.GetJSON)
	r.Post("/update/", metricHendlers.UpdateJSON)

	return r
}

func responseWithError(w http.ResponseWriter, code int, err error) {
	resp := fmt.Sprintf("%d: %s", code, err.Error())
	logger.Log.Error(resp)
	http.Error(w, resp, code)
}

func responseWithCode(w http.ResponseWriter, code int) {
	resp := fmt.Sprintf("%d: %s", code, http.StatusText(code))
	logger.Log.Debug(resp)
	w.WriteHeader(code)
}
