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
)

type metricSaver interface {
	Save(metric, metricType, value string) error
}

type updateHandler struct {
	saver metricSaver
}

func NewUpdateHandler(saver metricSaver) updateHandler {
	return updateHandler{saver}
}

func (h updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	metricType := parts[0]
	metric := parts[1]
	value := parts[2]

	err := h.saver.Save(metric, metricType, value)
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
