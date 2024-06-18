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
	// metricService provides methods for handling metrics.
	metricService interface {
		// Push pushes a metric with the given name, kind, and value.
		Push(ctx context.Context, name, kind, value string) error

		// PushCounter pushes a counter metric with the given name and value.
		PushCounter(ctx context.Context, name string, value metric.Counter) (metric.Counter, error)

		// PushGauge pushes a gauge metric with the given name and value.
		PushGauge(ctx context.Context, name string, value metric.Gauge) (metric.Gauge, error)

		// PushBatch pushes a batch of records to the storage.
		PushBatch(ctx context.Context, records []storage.Record) error

		// Get retrieves a record with the given name and kind.
		Get(ctx context.Context, name, kind string) (*storage.Record, error)

		// GetAll retrieves all records from the storage.
		GetAll(ctx context.Context) []storage.Record

		// Ping pings the service.
		Ping(ctx context.Context) error
	}

	// MetricHandlers contains methods for handling metrics.
	MetricHandlers struct {
		// service is the metric service.
		service metricService
		// logger is the logger for the metric handlers.
		logger *zap.Logger
	}
)

// NewMetricHandlers creates a new MetricHandlers instance.
func newMetricHandlers(s metricService, logger *zap.Logger) MetricHandlers {
	return MetricHandlers{
		service: s,
		logger:  logger,
	}
}

// List handles the HTTP GET request to retrieve all records.
// It sets the Content-Type header to "text/html; charset=UTF-8".
// It gets all records from the service and writes them to the response writer.
// It then responds with status code http.StatusOK.

// List retrieves all records and writes them to the ResponseWriter in a tabular format.
//
// Parameters:
//   - w: The ResponseWriter to write the response.
//   - r: The Request object representing the HTTP request.
//
// Returns: None.
func (h MetricHandlers) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	records := h.service.GetAll(r.Context())
	for _, v := range records {
		io.WriteString(w, fmt.Sprintf("%s\t%s\n", v.GetName(), v.GetValue().String()))
	}

	responseWithCode(w, http.StatusOK, h.logger)
}

// Get retrieves a specific record based on the provided kind and name.
//
// Parameters:
//   - w: The ResponseWriter to write the response.
//   - r: The Request object representing the HTTP request.
//
// Returns: None.
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

// Update updates a specific record with the provided kind, name, and value.
//
// Parameters:
//   - w: The ResponseWriter to write the response.
//   - r: The Request object representing the HTTP request.
//
// Returns: None.
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

//	Ping godoc
// 
//	@Summary Ping
//	@Description Performs a health check by pinging the service.
//	@Tags ping
//	@ID ping
//	@Success 200
//	@Failure 500
//	@Router /ping [get]
//line for correct view in godoc
// Ping handles the HTTP GET request to ping the service.
// It responds with status code http.StatusOK if the service is healthy.
// It responds with status code http.StatusInternalServerError if there is an error pinging the service.
func (h MetricHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(r.Context()); err != nil {
		responseWithCode(w, http.StatusInternalServerError, h.logger)
		return
	}

	responseWithCode(w, http.StatusOK, h.logger)
}
