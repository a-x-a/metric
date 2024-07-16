// Package encoder описывает middleware для работа со сжатыми HTTP запросами и ответами.
package encoder

import (
	"compress/gzip"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// DecompressMiddleware middleware для распаковки данных.
func DecompressMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encoding := r.Header.Get("Content-Encoding")
			if len(encoding) == 0 {
				logger.Info("got uncompressed request", zap.String("encoding", encoding))
				next.ServeHTTP(w, r)
				return
			}

			if !strings.Contains(encoding, "gzip") {
				logger.Info("compressed method not supported", zap.String("method", "gzip"))
				next.ServeHTTP(w, r)
				return
			}

			logger.Info("request compressed", zap.String("method", encoding))

			cr, err := newCompressReader(r.Body)
			if err != nil {
				logger.Error("compress reader", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}

			defer cr.Close()

			r.Body = cr

			next.ServeHTTP(w, r)
		})
	}
}

// CompressMiddleware middleware для упаковки данных.
func CompressMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				logger.Info("compression not supported by client", zap.String("method", "gzip"))
				next.ServeHTTP(w, r)
				return
			}

			logger.Info("compression supported by client", zap.String("method", "gzip"))

			zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				logger.Error("compress writer", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}

			defer zw.Close()

			w.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: zw}, r)
		})
	}
}
