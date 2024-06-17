package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/encoder"
	"github.com/a-x-a/go-metric/internal/logger"
	"github.com/a-x-a/go-metric/internal/signer"
)

func NewRouter(s metricService, log *zap.Logger, key string) http.Handler {
	metricHendlers := newMetricHandlers(s, log)

	r := chi.NewRouter()

	r.Use(logger.LoggerMiddleware(log))
	r.Use(encoder.DecompressMiddleware(log))
	r.Use(encoder.CompressMiddleware(log))

	r.Mount("/debug", middleware.Profiler())

	r.Get("/", metricHendlers.List)

	r.Post("/value", metricHendlers.GetJSON)
	r.Post("/value/", metricHendlers.GetJSON)
	r.Get("/value/{kind}/{name}", metricHendlers.Get)

	r.Post("/update", metricHendlers.UpdateJSON)
	r.Post("/update/", metricHendlers.UpdateJSON)
	r.Post("/update/{kind}/{name}/{value}", metricHendlers.Update)

	updateGroup := r.Group(nil)
	if len(key) != 0 {
		updateGroup.Use(signer.SignerMiddleware(log, key))
	}
	updateGroup.Post("/updates", metricHendlers.UpdateBatch)
	updateGroup.Post("/updates/", metricHendlers.UpdateBatch)

	r.Get("/ping", metricHendlers.Ping)
	r.Get("/ping/", metricHendlers.Ping)

	return r
}

func responseWithError(w http.ResponseWriter, code int, err error, logger *zap.Logger) {
	resp := fmt.Sprintf("%d: %s", code, err.Error())
	logger.Error(resp)
	http.Error(w, resp, code)
}

func responseWithCode(w http.ResponseWriter, code int, logger *zap.Logger) {
	resp := fmt.Sprintf("%d: %s", code, http.StatusText(code))
	logger.Debug(resp)
	w.WriteHeader(code)
}
