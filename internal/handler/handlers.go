package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/go-chi/chi"
)

type metricHandlers struct {
	metricService metricservice.MetricService
}

func newMetricHandlers(metricService metricservice.MetricService) metricHandlers {
	return metricHandlers{metricService}
}

func (h metricHandlers) List(w http.ResponseWriter, r *http.Request) {
	records := h.metricService.GetAll()

	for _, v := range records {
		io.WriteString(w, fmt.Sprintf("%s\t%s\n", v.GetName(), v.GetValue().String()))
	}

	ok(w)
}

func (h metricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")

	value, err := h.metricService.Get(name, kind)
	if err != nil {
		notFound(w)
		return
	}

	w.Write([]byte(value))
	ok(w)
}

func (h metricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	err := h.metricService.Push(name, kind, value)
	if err != nil {
		badRequest(w)
		return
	}

	ok(w)
}
