package encoder

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/logger"
)

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера.
// декомпрессировать получаемые от клиента данные.

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
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
