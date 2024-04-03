package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-x-a/go-metric/internal/adapter"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

func (h metricHandlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	data := &adapter.RequestMetric{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	switch data.MType {
	case "counter":
		val, err := h.service.PushCounter(data.ID, metric.Counter(*data.Delta))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err)
			return
		}

		newDelta := int64(val)
		data.Delta = &newDelta

	case "gauge":
		val, err := h.service.PushGauge(data.ID, metric.Gauge(*data.Value))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err)
			return
		}

		newValue := float64(val)
		data.Value = &newValue

		fmt.Println("data:=", data.ID, data.MType, *data.Value)
	default:
		responseWithCode(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	responseWithCode(w, http.StatusOK)
}

func (h metricHandlers) GetJSON(w http.ResponseWriter, r *http.Request) {
	data := &adapter.RequestMetric{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	value, err := h.service.Get(data.ID, data.MType)
	if err != nil {
		responseWithCode(w, http.StatusNotFound)
		return
	}

	switch data.MType {
	case "counter":
		val, err := metric.ToCounter(value)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err)
			return
		}

		newDelta := int64(val)
		data.Delta = &newDelta

	case "gauge":
		val, err := metric.ToGauge(value)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err)
			return
		}

		newValue := float64(val)
		data.Value = &newValue

	default:
		responseWithCode(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	responseWithCode(w, http.StatusOK)
}
