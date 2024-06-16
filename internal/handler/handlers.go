package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	metricService interface {
		Push(ctx context.Context, name, kind, value string) error
		PushCounter(ctx context.Context, name string, value metric.Counter) (metric.Counter, error)
		PushGauge(ctx context.Context, name string, value metric.Gauge) (metric.Gauge, error)
		PushBatch(ctx context.Context, records []storage.Record) error
		Get(ctx context.Context, name, kind string) (*storage.Record, error)
		GetAll(ctx context.Context) []storage.Record
		Ping(ctx context.Context) error
	}

	MetricHandlers struct {
		service metricService
		logger  *zap.Logger
	}
)

func newMetricHandlers(s metricService, logger *zap.Logger) MetricHandlers {
	return MetricHandlers{
		service: s,
		logger:  logger,
	}
}

func (h MetricHandlers) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	records := h.service.GetAll(r.Context())
	for _, v := range records {
		io.WriteString(w, fmt.Sprintf("%s\t%s\n", v.GetName(), v.GetValue().String()))
	}

	responseWithCode(w, http.StatusOK, h.logger)
}

func (h MetricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")

	record, err := h.service.Get(r.Context(), name, kind)
	if err != nil {
		responseWithCode(w, http.StatusNotFound, h.logger)
		return
	}

	value := record.GetValue().String()
	w.Write([]byte(value))

	responseWithCode(w, http.StatusOK, h.logger)
}

func (h MetricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	err := h.service.Push(r.Context(), name, kind, value)
	if err != nil {
		responseWithCode(w, http.StatusBadRequest, h.logger)
		return
	}

	responseWithCode(w, http.StatusOK, h.logger)
}

func (h MetricHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(r.Context()); err != nil {
		responseWithCode(w, http.StatusInternalServerError, h.logger)
		return
	}

	responseWithCode(w, http.StatusOK, h.logger)
}
