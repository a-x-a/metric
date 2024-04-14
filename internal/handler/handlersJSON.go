package handler

import (
	"encoding/json"
	"net/http"

	"github.com/a-x-a/go-metric/internal/adapter"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

func (h metricHandlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
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
		val, err := h.service.PushCounter(data.ID, metric.Counter(*data.Delta))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err, h.logger)
			return
		}

		newDelta := int64(val)
		data.Delta = &newDelta

	case metric.KindGauge:
		val, err := h.service.PushGauge(data.ID, metric.Gauge(*data.Value))
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

func (h metricHandlers) GetJSON(w http.ResponseWriter, r *http.Request) {
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

	value, err := h.service.Get(data.ID, data.MType)
	if err != nil {
		responseWithCode(w, http.StatusNotFound, h.logger)
		return
	}

	switch kind {
	case metric.KindCounter:
		val, err := metric.ToCounter(value)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err, h.logger)
			return
		}

		newDelta := int64(val)
		data.Delta = &newDelta

	case metric.KindGauge:
		val, err := metric.ToGauge(value)
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
