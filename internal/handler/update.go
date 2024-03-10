/*
  - Принимать метрики по протоколу HTTP методом `POST`.
  - Принимать данные в формате:
    `http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>`, `Content-Type: text/plain`.
  - При успешном приёме возвращать `http.StatusOK`.
  - При попытке передать запрос без имени метрики возвращать `http.StatusNotFound`.
  - При попытке передать запрос с некорректным типом метрики или значением возвращать `http.StatusBadRequest`.
*/
package handler

import (
	"net/http"
	"strings"

	"github.com/a-x-a/go-metric/internal/service/metricservice"
)

type metricHandlers struct {
	metricService metricservice.MetricService
}

func newMetricHandlers(metricService metricservice.MetricService) metricHandlers {
	return metricHandlers{metricService}
}

func (h metricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	// принимаем метрики методом POST
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/update/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 {
		notFound(w)
		return
	}

	kind := parts[0]
	name := parts[1]
	value := parts[2]

	err := h.metricService.Push(name, kind, value)
	if err != nil {
		badRequest(w)
		return
	}

	ok(w)
}

func ok(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

func badRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
}
