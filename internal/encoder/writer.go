package encoder

import (
	"compress/gzip"
	"net/http"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/logger"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера.
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки.
type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	// zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	// if err != nil {
	// 	return nil, err
	// }

	return &compressWriter{
		ResponseWriter: w,
		// zw:             zw,
	}
}

// func (c *compressWriter) Header() http.Header {
// 	return c.Header()
// }

func (c *compressWriter) Write(p []byte) (int, error) {
	contentType := c.Header().Get("Content-Type")
	if !isSupportedContentType(contentType) {
		logger.Log.Debug("сжатие не поддерживается", zap.String("ContentType", contentType))
		return c.ResponseWriter.Write(p)
	}

	if c.zw == nil {
		zw, err := gzip.NewWriterLevel(c.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			logger.Log.Error("compressWriter", zap.Error(err))
			return 0, err //c.ResponseWriter.Write(p)
		}

		c.zw = zw
	}

	c.Header().Set("Content-Encoding", "gzip")

	return c.zw.Write(p)
}

// func (c *compressWriter) WriteHeader(statusCode int) {
// 	if statusCode < http.StatusMultipleChoices {
// 		c.ResponseWriter.Header().Set("Content-Encoding", "gzip")
// 	}
// 	c.ResponseWriter.WriteHeader(statusCode)
// }

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() {
	if c.zw != nil {
		c.zw.Close()
	}
}
