package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Router(s metricService) http.Handler {
	metricHendlers := newMetricHandlers(s)
	r := chi.NewRouter()
	// r.Use(logging.RequestsLogger)
	r.Get("/", metricHendlers.List)
	r.Get("/value/{kind}/{name}", metricHendlers.Get)
	r.Post("/update/{kind}/{name}/{value}", metricHendlers.Update)

	return r
}

func ok(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

func badRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
}

// func internalServerError(w http.ResponseWriter) {
// 	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
// 	w.WriteHeader(http.StatusInternalServerError)
// }
