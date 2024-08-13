// Package sender отвечает за оправку данных от агента к серверу.
package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/security"
	"github.com/a-x-a/go-metric/pkg/grpcapi"
)

type GRPCSender struct {
	baseURL   string
	client    *grpc.ClientConn
	signer    *security.Signer
	cryptoKey security.PublicKey
	batch     chan *grpcapi.Metric
	err       error
}

func NewGRPCSender(serverAddress string, secret string, publicKey security.PublicKey) *GRPCSender {
	sgnr := security.NewSigner(secret)

	return &GRPCSender{
		baseURL:   serverAddress,
		client:    nil,
		signer:    sgnr,
		cryptoKey: publicKey,
		batch:     make(chan *grpcapi.Metric, 1024),
		err:       nil,
	}
}

func (gs *GRPCSender) doSend(ctx context.Context, batch []*grpcapi.Metric) error {
	var err error
	if len(batch) == 0 {
		return fmt.Errorf("metrics send: batch is empty")
	}

	data, err := json.Marshal(batch)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("metrics send: data is empty")
	}

	if gs.client == nil {
		gs.client, err = grpc.DialContext(
			ctx,
			gs.baseURL,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}
	}

	// var buf bytes.Buffer
	// if err = encoder.Encoding(data, &buf); err != nil {
	// 	return err
	// }

	// if hs.cryptoKey != nil {
	// 	b, err := security.Encrypt(&buf, hs.cryptoKey)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	buf = *b
	// }

	md := make(map[string]string)
	ip, err := getOutboundIP()
	if err != nil {
		return err
	}

	md["X-Real-IP"] = ip.String()

	// if gs.signer != nil {
	// 	var hash []byte
	// 	hash, err = gs.signer.Hash(data)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	md["HashSHA256"] = hex.EncodeToString(hash)
	// }

	req := &grpcapi.UpdateBatchMetricRequestV1{Data: batch}
	grpcClient := grpcapi.NewMetricsClient(gs.client)
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(md))

	_, err = grpcClient.UpdateBatch(ctx, req)

	if err != nil {
		return fmt.Errorf("metrics send failed: (%w)", err)
	}

	return nil
}

func (gs *GRPCSender) worker(ctx context.Context) {
	data := make([]*grpcapi.Metric, 0, len(gs.batch))

	for r := range gs.batch {
		data = append(data, r)
	}

	gs.err = gs.doSend(ctx, data)
}

func (gs *GRPCSender) Add(name string, value metric.Metric) Sender {
	if gs.err != nil {
		return gs
	}

	val := grpcapi.Metric{
		Id:    name,
		Mtype: value.Kind(),
	}

	switch {
	case value.IsCounter():
		v, ok := value.(metric.Counter)
		if !ok {
			gs.err = fmt.Errorf("fail to convert counter %v", value)
			return gs
		}
		val.Value = &grpcapi.Metric_Counter{
			Counter: int64(v),
		}
	case value.IsGauge():
		v, ok := value.(metric.Gauge)
		if !ok {
			gs.err = fmt.Errorf("fail to convert gauge %v", value)
			return gs
		}
		val.Value = &grpcapi.Metric_Gauge{
			Gauge: float64(v),
		}
	}

	gs.batch <- &val

	return gs
}

func (gs *GRPCSender) Send(ctx context.Context, rateLimit int) error {
	wg := sync.WaitGroup{}
	wg.Add(rateLimit)
	for i := 0; i < rateLimit; i++ {
		go func() {
			defer wg.Done()
			gs.worker(ctx)
		}()
	}

	close(gs.batch)

	wg.Wait()

	return gs.err
}

func (gs *GRPCSender) Reset() {
	gs.batch = make(chan *grpcapi.Metric, 1024)
	gs.err = nil
}

func (gs *GRPCSender) Close() error {
	if gs.client == nil {
		return nil
	}

	return gs.client.Close()
}
