package encoder

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/logger"
)

func isSupportedContentType(contentType string) bool {
	supportedContentType :=
		[...]string{"application/json", "text/html; charset=utf-8"}

	for _, v := range supportedContentType {
		if strings.Contains(v, contentType) {
			return true
		}
	}

	return false
}

func DecompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Content-Encoding")
		if len(encoding) == 0 {
			logger.Log.Info("got uncompressed request", zap.String("encoding", encoding))
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(encoding, "gzip") {
			logger.Log.Info("compressed method not supported", zap.String("method", "gzip"))
			next.ServeHTTP(w, r)
			return
		}

		logger.Log.Info("request compressed", zap.String("method", encoding))

		cr, err := newCompressReader(r.Body)
		if err != nil {
			logger.Log.Error("compress reader", zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		defer cr.Close()

		r.Body = cr

		next.ServeHTTP(w, r)
	})
}

func CompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			logger.Log.Info("compression not supported by client", zap.String("method", "gzip"))
			next.ServeHTTP(w, r)
			return
		}

		logger.Log.Info("compression supported by client", zap.String("method", "gzip"))

		cw, err := newCompressWriter(w)
		if err != nil {
			logger.Log.Error("compress writer", zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		defer cw.Close()

		next.ServeHTTP(cw, r)
	})
}
