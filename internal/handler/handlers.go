package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	metricService interface {
		Push(name, kind, value string) error
		PushCounter(name string, value metric.Counter) (metric.Counter, error)
		PushGauge(name string, value metric.Gauge) (metric.Gauge, error)
		Get(name, kind string) (string, error)
		GetAll() []storage.Record
	}
	metricHandlers struct {
		service metricService
	}
)

func newMetricHandlers(s metricService) metricHandlers {
	return metricHandlers{s}
}

func (h metricHandlers) List(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	records := h.service.GetAll()
	for _, v := range records {
		io.WriteString(w, fmt.Sprintf("%s\t%s\n", v.GetName(), v.GetValue().String()))
	}

	responseWithCode(w, http.StatusOK)
}

func (h metricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")

	value, err := h.service.Get(name, kind)
	if err != nil {
		responseWithCode(w, http.StatusNotFound)
		return
	}

	w.Write([]byte(value))

	responseWithCode(w, http.StatusOK)
}

func (h metricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	err := h.service.Push(name, kind, value)
	if err != nil {
		responseWithCode(w, http.StatusBadRequest)
		return
	}

	responseWithCode(w, http.StatusOK)
}
