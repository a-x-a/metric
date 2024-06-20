package handler

import (
	"encoding/json"
	"net/http"

	"github.com/a-x-a/go-metric/internal/adapter"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

//	UpdateBatch godoc
//
//	@Summary		UpdateBatch
//	@Description	Обновляет текущие значения метрик из набора.
//	@Tags			updatebatch
//	@ID				updatebatch
//	@Produce		json
//	@Param			data	body	[]adapter.RequestMetric	true
//	@Failure		404
//	@Failure		500
//	@Router			/updates [post]
//
// Update обновляет значение метрики с указанным именем и типом.
// В случае ошибки, статус ответа 404
func (h MetricHandlers) UpdateBatch(w http.ResponseWriter, r *http.Request) {
	data := make([]adapter.RequestMetric, 0)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		responseWithError(w, http.StatusBadRequest, err, h.logger)
		return
	}

	records := make([]storage.Record, 0)

	for _, v := range data {
		record, err := storage.NewRecord(v.ID)
		if err != nil {
			responseWithCode(w, http.StatusBadRequest, h.logger)
			return
		}

		kind, err := metric.GetKind(v.MType)
		if err != nil {
			responseWithCode(w, http.StatusBadRequest, h.logger)
			return
		}

		switch kind {
		case metric.KindCounter:
			if v.Delta == nil {
				responseWithError(w, http.StatusBadRequest, err, h.logger)
				return
			}
			val := metric.Counter(*v.Delta)
			record.SetValue(val)
		case metric.KindGauge:
			if v.Value == nil {
				responseWithError(w, http.StatusBadRequest, err, h.logger)
				return
			}
			val := metric.Gauge(*v.Value)
			record.SetValue(val)
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		responseWithCode(w, http.StatusBadRequest, h.logger)
		return
	}

	if err := h.service.PushBatch(r.Context(), records); err != nil {
		responseWithError(w, http.StatusInternalServerError, err, h.logger)
		return
	}
}
