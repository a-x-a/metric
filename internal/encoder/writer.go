package encoder

import (
	"compress/gzip"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/logger"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (c compressWriter) Write(p []byte) (int, error) {
	contentType := c.Header().Get("Content-Type")
	if !isSupportedContentType(contentType) {
		logger.Log.Debug("сжатие не поддерживается", zap.String("ContentType", contentType))
		return c.ResponseWriter.Write(p)
	}

	if c.Writer == nil {
		zw, err := gzip.NewWriterLevel(c.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			logger.Log.Error("compressWriter", zap.Error(err))
			return 0, err //c.ResponseWriter.Write(p)
		}

		c.Writer = zw
	}

	// c.Header().Set("Content-Encoding", "gzip")

	return c.Writer.Write(p)
}
