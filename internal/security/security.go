package security

import (
	"bytes"
	"encoding/pem"
	"errors"
	"io"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
)

var (
	ErrNotSupportedFormatKey = errors.New("не поддерживаемый тип ключа")
	ErrNotPEMFormatFile      = errors.New("указанный файл не содержит ключ в формате PEM")
	ErrUntrustedSource       = errors.New("запрос получен из не доверенного истоочника")
)

// getRawKey считывает ключ из PEM файла.
func getRawKey(path string) (*pem.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rawKey, _ := pem.Decode(data)
	if rawKey == nil {
		return nil, ErrNotPEMFormatFile
	}

	return rawKey, nil
}

// DecryptMiddleware HTTP middleware расшифровывает полученный запрос с использованием алгоритма RSA.
func DecryptMiddleware(logger *zap.Logger, key PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			msg, err := Decrypt(r.Body, key)
			if err != nil {
				logger.Error("security.DecryptMiddleware Decrypt", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(msg.Bytes()))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// TrustedSubnetMiddleware HTTP middleware выполняет проверку на разрешённые IP адреса.
func TrustedSubnetMiddleware(logger *zap.Logger, subnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			xrip := r.Header.Get("X-Real-IP")
			rip := net.ParseIP(xrip)
			if !subnet.Contains(rip) {
				logger.Error("security.TrustedSubnetMiddleware subnet.Contains", zap.Error(ErrUntrustedSource))
				http.Error(w, "", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
