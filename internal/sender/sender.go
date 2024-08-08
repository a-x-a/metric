// Package sender отвечает за оправку данных от агента к серверу.
package sender

import (
	"context"
	"time"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/security"
)

type Sender interface {
	Add(name string, value metric.Metric) Sender
	Send(ctx context.Context, rateLimit int) error
	Close() error
	Reset()
}

const (
	TransportHTTP = "http"
	TransportGRPC = "grpc"
)

func New(transport string, address string, pollInterval time.Duration, secret string, publicKey security.PublicKey) Sender {
	if transport == TransportGRPC {
		return NewGRPCSender(address, secret, publicKey)
	}

	return NewHTTPSender(address, pollInterval, secret, publicKey)
}
