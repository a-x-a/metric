package sender

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/a-x-a/go-metric/internal/encoder"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/security"
)

type HttpSender struct {
	baseURL   string
	client    *http.Client
	signer    *security.Signer
	cryptoKey security.PublicKey
	batch     chan metric.RequestMetric
	err       error
}

func NewHTTPSender(serverAddress string, timeout time.Duration, secret string, publicKey security.PublicKey) *HttpSender {
	baseURL := fmt.Sprintf("http://%s", serverAddress)
	client := &http.Client{Timeout: timeout}
	sgnr := security.NewSigner(secret)

	return &HttpSender{
		baseURL:   baseURL,
		client:    client,
		signer:    sgnr,
		cryptoKey: publicKey,
		batch:     make(chan metric.RequestMetric, 1024),
		err:       nil,
	}
}

func (hs *HttpSender) doSend(ctx context.Context, batch []metric.RequestMetric) error {
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

	var buf bytes.Buffer
	if err = encoder.Encoding(data, &buf); err != nil {
		return err
	}

	if hs.cryptoKey != nil {
		b, err := security.Encrypt(&buf, hs.cryptoKey)
		if err != nil {
			return err
		}
		buf = *b
	}

	ip, err := getOutboundIP()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hs.baseURL+"/updates", &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("X-Real-IP", ip.String())

	if hs.signer != nil {
		var hash []byte
		hash, err = hs.signer.Hash(data)
		if err != nil {
			return err
		}
		req.Header.Set("HashSHA256", hex.EncodeToString(hash))
	}

	resp, err := hs.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if _, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("metrics send failed: (%d)", resp.StatusCode)
	}

	return nil
}

func (hs *HttpSender) worker(ctx context.Context) {
	data := make([]metric.RequestMetric, 0, len(hs.batch))

	for r := range hs.batch {
		data = append(data, r)
	}

	hs.err = hs.doSend(ctx, data)
}

func (hs *HttpSender) Add(name string, value metric.Metric) Sender {
	if hs.err != nil {
		return hs
	}

	val := metric.RequestMetric{
		ID:    name,
		MType: value.Kind(),
	}

	switch {
	case value.IsCounter():
		v, ok := value.(metric.Counter)
		if !ok {
			hs.err = fmt.Errorf("fail to convert counter %v", value)
			return hs
		}
		val.Delta = (*int64)(&v)
	case value.IsGauge():
		v, ok := value.(metric.Gauge)
		if !ok {
			hs.err = fmt.Errorf("fail to convert gauge %v", value)
			return hs
		}
		val.Value = (*float64)(&v)
	}

	hs.batch <- val

	return hs
}

func (hs *HttpSender) Send(ctx context.Context, rateLimit int) error {
	wg := sync.WaitGroup{}
	wg.Add(rateLimit)
	for i := 0; i < rateLimit; i++ {
		go func() {
			defer wg.Done()
			hs.worker(ctx)
		}()
	}

	close(hs.batch)

	wg.Wait()

	return hs.err
}

func (hs *HttpSender) Reset() {
	hs.batch = make(chan metric.RequestMetric, 1024)
	hs.err = nil
}

func (hs *HttpSender) Close() error {
	return nil
}
