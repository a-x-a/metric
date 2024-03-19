package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	metricService interface {
		Push(name, kind, value string) error
		Get(name, kind string) (string, error)
		GetAll() []storage.Record
	}
)
type metricHandlers struct {
	service metricService
}

func newMetricHandlers(s metricService) metricHandlers {
	return metricHandlers{s}
}

func (h metricHandlers) List(w http.ResponseWriter, r *http.Request) {
	records := h.service.GetAll()
	for _, v := range records {
		io.WriteString(w, fmt.Sprintf("%s\t%s\n", v.GetName(), v.GetValue().String()))
	}
	ok(w)
}

func (h metricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")
	value, err := h.service.Get(name, kind)
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
	err := h.service.Push(name, kind, value)
	if err != nil {
		badRequest(w)
		return
	}
	ok(w)
}
