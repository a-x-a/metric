package middlewarewithlogger

import (
	"compress/gzip"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func (m middlewareWithLogger) Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Content-Encoding")
		if len(encoding) == 0 {
			m.logger.Info("got uncompressed request", zap.String("encoding", encoding))
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(encoding, "gzip") {
			m.logger.Info("compressed method not supported", zap.String("method", "gzip"))
			next.ServeHTTP(w, r)
			return
		}

		m.logger.Info("request compressed", zap.String("method", encoding))

		cr, err := newCompressReader(r.Body)
		if err != nil {
			m.logger.Error("compress reader", zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		defer cr.Close()

		r.Body = cr

		next.ServeHTTP(w, r)
	})
}

func (m middlewareWithLogger) Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			m.logger.Info("compression not supported by client", zap.String("method", "gzip"))
			next.ServeHTTP(w, r)
			return
		}

		m.logger.Info("compression supported by client", zap.String("method", "gzip"))

		zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			m.logger.Error("compress writer", zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		defer zw.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: zw}, r)
	})
}
