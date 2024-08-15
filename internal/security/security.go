package security

import (
	"bytes"
	"context"
	"encoding/pem"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

// UnaryRequestsInterceptor GRPC interceptor выполняет проверку на разрешённые IP адреса.
func UnaryRequestsInterceptor(logger *zap.Logger, subnet *net.IPNet) grpc.UnaryServerInterceptor {
	fn := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger.Info("grpc request",
			zap.Any("srv", info.Server),
			zap.String("method", info.FullMethod),
		)

		if !strings.HasSuffix(info.FullMethod, "UpdateBatch") {
			return handler(ctx, req)
		}

		var rip net.IP

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get("x-real-ip")
			if len(values) > 0 {
				rip = net.ParseIP(values[0])
			}
		}

		logger.Info("client address", zap.String("x-real-ip", rip.String()))

		if !subnet.Contains(rip) {
			logger.Error("security.UnaryRequestsInterceptor subnet.Contains", zap.Error(ErrUntrustedSource))
			return nil, status.Error(codes.PermissionDenied, ErrUntrustedSource.Error())
		}

		return handler(ctx, req)
	}

	return fn
}
