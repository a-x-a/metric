package handler

import (
	"encoding/json"
	"net/http"

	"github.com/a-x-a/go-metric/internal/adapter"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

//	UpdateJSON godoc
//
//	@Summary		UpdateJSON
//	@Description	Обновляет текущее значение метрики с указанным имененм и типом.
//	@Tags			update
//	@ID				updateJSON
//	@Produce		json
//	@Param			data	body	adapter.RequestMetric	true	"Параметры метрики: имя, тип, значение"
//	@Success		200
//	@Failure		404
//	@Failure		500
//	@Router			/update [post]
//
// UpdateJSON обновляет текущее значение метрики с указанным именем и типом полученные в формате JSON.
//
//line for correct view in godoc.
func (h MetricHandlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	data := &adapter.RequestMetric{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		responseWithError(w, http.StatusBadRequest, err, h.logger)
		return
	}

	kind, err := metric.GetKind(data.MType)
	if err != nil {
		responseWithCode(w, http.StatusBadRequest, h.logger)
		return
	}

	switch kind {
	case metric.KindCounter:
		val, err := h.service.PushCounter(r.Context(), data.ID, metric.Counter(*data.Delta))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err, h.logger)
			return
		}

		newDelta := int64(val)
		data.Delta = &newDelta

	case metric.KindGauge:
		val, err := h.service.PushGauge(r.Context(), data.ID, metric.Gauge(*data.Value))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err, h.logger)
			return
		}

		newValue := float64(val)
		data.Value = &newValue
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		responseWithError(w, http.StatusInternalServerError, err, h.logger)
		return
	}

	responseWithCode(w, http.StatusOK, h.logger)
}

//	GetJSON godoc
//
//	@Summary		GetJSON
//	@Description	Возвращает текущее значение метрики в формате JSON с указанным имененм и типом.
//	@Tags			value
//	@ID				getJSON
//	@Accept			json
//	@Produce		json
//	@Param			data	body	adapter.RequestMetric	true	"Параметры метрики: имя, тип"
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/value [get]
//
// GetJSON возвращает текущее значение метрики в формате JSON с указанным имененм и типом в формате.
//
//line for correct view in godoc.
func (h MetricHandlers) GetJSON(w http.ResponseWriter, r *http.Request) {
	data := &adapter.RequestMetric{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		responseWithError(w, http.StatusBadRequest, err, h.logger)
		return
	}

	kind, err := metric.GetKind(data.MType)
	if err != nil {
		responseWithCode(w, http.StatusBadRequest, h.logger)
		return
	}

	record, err := h.service.Get(r.Context(), data.ID, data.MType)
	if err != nil {
		responseWithCode(w, http.StatusNotFound, h.logger)
		return
	}

	value := record.GetValue()

	switch kind {
	case metric.KindCounter:
		val, ok := value.(metric.Counter)
		if !ok {
			responseWithError(w, http.StatusInternalServerError, err, h.logger)
			return
		}

		newDelta := int64(val)
		data.Delta = &newDelta

	case metric.KindGauge:
		val, ok := value.(metric.Gauge)
		if !ok {
			responseWithError(w, http.StatusInternalServerError, err, h.logger)
			return
		}

		newValue := float64(val)
		data.Value = &newValue
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		responseWithError(w, http.StatusInternalServerError, err, h.logger)
		return
	}

	responseWithCode(w, http.StatusOK, h.logger)
}
