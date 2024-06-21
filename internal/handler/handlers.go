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
	// metricService содержит описание методов для обработки запросов.
	metricService interface {
		// Push добавляет метрику с указанным именем, типом и значением.
		Push(ctx context.Context, name, kind, value string) error

		// PushCounter добавляет метрику с указанным именем с типом counter и значением.
		PushCounter(ctx context.Context, name string, value metric.Counter) (metric.Counter, error)

		// PushGauge добавляет метрику с указанным именем с типом gauge и значением.
		PushGauge(ctx context.Context, name string, value metric.Gauge) (metric.Gauge, error)

		// PushBatch добавляет набор метрик.
		PushBatch(ctx context.Context, records []storage.Record) error

		// Get получает текущее значение метрики с указанным именем и типом.
		Get(ctx context.Context, name, kind string) (*storage.Record, error)

		// GetAll получает текущее значение всех метрик.
		GetAll(ctx context.Context) []storage.Record

		// Ping проверяет состояние сервиса.
		Ping(ctx context.Context) error
	}

	// MetricHandlers содержит методы для обработки запросов.
	MetricHandlers struct {
		service metricService
		logger  *zap.Logger
	}
)

// newMetricHandlers создаёт новый экземпляр объекта MetricHandlers.
func newMetricHandlers(s metricService, logger *zap.Logger) MetricHandlers {
	return MetricHandlers{
		service: s,
		logger:  logger,
	}
}

//	List godoc
//
//	@Summary		List
//	@Description	Возвращает список полученных от сервиса метрик ввиде обычного текста.
//	@Tags			list
//	@ID				list
//	@Produce		html
//	@Success		200
//	@Router			/list [get]
//
// List обрабатывает HTTP-GET-запрос на получение списка текущих значений всех метрик.
// Возвращает список полученных от сервиса метрик ввиде обычного текста.
func (h MetricHandlers) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	records := h.service.GetAll(r.Context())
	for _, v := range records {
		io.WriteString(w, fmt.Sprintf("%s\t%s\n", v.GetName(), v.GetValue().String()))
	}

	responseWithCode(w, http.StatusOK, h.logger)
}

//	Get godoc
//
//	@Summary		Get
//	@Description	Возвращает текущее значение метрики с указанным имененм и типом.
//	@Tags			value
//	@ID				get
//	@Produce		html
//	@Param			kind	path	string	true	"Тип метрики"
//	@Param			name	path	string	true	"Имя метрики"
//	@Success		200
//	@Failure		404
//	@Router			/value/{kind}/{name} [get]
//
// Get возвращает текущее значение метрики с указанным имененм и типом.
// В случае ошибки, статус ответа 404
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

//	Update godoc
//
//	@Summary		Update
//	@Description	Обновляет текущее значение метрики с указанным имененм и типом.
//	@Tags			update
//	@ID				update
//	@Produce		html
//	@Param			kind	path	string	true	"Тип метрики"
//	@Param			name	path	string	true	"Имяметрики"
//	@Param			value	path	string	true	"Значение метрики"
//	@Success		200
//	@Failure		404
//	@Router			/update/{kind}/{name}/{value} [post]
//
// Update обновляет значение метрики с указанным именем и типом.
// В случае ошибки, статус ответа 404
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
//	@Summary		Ping
//	@Description	Performs a health check by pinging the service.
//	@Tags			ping
//	@ID				ping
//	@Success		200
//	@Failure		500
//	@Router			/ping [get]
//
// Ping проверяет состояние сервиса.
// В случае ошибки, статус ответа 500
func (h MetricHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(r.Context()); err != nil {
		responseWithCode(w, http.StatusInternalServerError, h.logger)
		return
	}

	responseWithCode(w, http.StatusOK, h.logger)
}
