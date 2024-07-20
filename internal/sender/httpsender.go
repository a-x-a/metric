// Package sender отвечает за оправку данных от агента к серверу.
package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-x-a/go-metric/internal/adapter"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/security"
)

type httpSender struct {
	baseURL   string
	client    *http.Client
	signer    *security.Signer
	cryptoKey security.PublicKey
	batch     chan adapter.RequestMetric
	err       error
}

func newHTTPSender(serverAddress string, timeout time.Duration, key string, cryptoKey security.PublicKey) httpSender {
	baseURL := fmt.Sprintf("http://%s", serverAddress)
	client := &http.Client{Timeout: timeout}
	sgnr := security.NewSigner(key)

	return httpSender{
		baseURL:   baseURL,
		client:    client,
		signer:    sgnr,
		cryptoKey: cryptoKey,
		batch:     make(chan adapter.RequestMetric, 1024),
		err:       nil,
	}
}

func (hs *httpSender) doSend(ctx context.Context, batch []adapter.RequestMetric) error {
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

	zw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return err
	}

	if _, err = zw.Write(data); err != nil {
		return err
	}

	if err = zw.Close(); err != nil {
		return err
	}

	if hs.cryptoKey != nil {
		b, err := security.Encrypt(&buf, hs.cryptoKey)
		if err != nil {
			return err
		}
		buf = *b
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hs.baseURL+"/updates", &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

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

func (hs *httpSender) add(rm adapter.RequestMetric) *httpSender {
	if hs.err != nil {
		return hs
	}

	hs.batch <- rm

	return hs
}

// SendMetrics отправляет метрики на сервер.
//
// Parameters:
// - ctx: контекст.
// - serverAddress: адрес сервера сбора метрик.
// - timeout: частота отправки метрик на сервер.
// - key: ключ подписи данных.
// - rateLimit: количество одновременно исходящих запросов на сервер.
// - stats: коллекция мсетрик для отправки.
// - сryptoKey публичныq ключ.
func SendMetrics(ctx context.Context, serverAddress string, timeout time.Duration, key string, rateLimit int, stats metric.Metrics, cryptoKey security.PublicKey) error {
	sender := newHTTPSender(serverAddress, timeout, key, cryptoKey)

	for i := 0; i < rateLimit; i++ {
		go sender.worker(ctx)
	}

	// отправляем метрики пакета runtime
	sender.
		add(adapter.NewUpdateRequestMetricGauge("Alloc", stats.Runtime.Alloc)).
		add(adapter.NewUpdateRequestMetricGauge("BuckHashSys", stats.Runtime.BuckHashSys)).
		add(adapter.NewUpdateRequestMetricGauge("Frees", stats.Runtime.Frees)).
		add(adapter.NewUpdateRequestMetricGauge("GCCPUFraction", stats.Runtime.GCCPUFraction)).
		add(adapter.NewUpdateRequestMetricGauge("GCSys", stats.Runtime.GCSys)).
		add(adapter.NewUpdateRequestMetricGauge("HeapAlloc", stats.Runtime.HeapAlloc)).
		add(adapter.NewUpdateRequestMetricGauge("HeapIdle", stats.Runtime.HeapIdle)).
		add(adapter.NewUpdateRequestMetricGauge("HeapInuse", stats.Runtime.HeapInuse)).
		add(adapter.NewUpdateRequestMetricGauge("HeapObjects", stats.Runtime.HeapObjects)).
		add(adapter.NewUpdateRequestMetricGauge("HeapReleased", stats.Runtime.HeapReleased)).
		add(adapter.NewUpdateRequestMetricGauge("HeapSys", stats.Runtime.HeapSys)).
		add(adapter.NewUpdateRequestMetricGauge("LastGC", stats.Runtime.LastGC)).
		add(adapter.NewUpdateRequestMetricGauge("Lookups", stats.Runtime.Lookups)).
		add(adapter.NewUpdateRequestMetricGauge("MCacheInuse", stats.Runtime.MCacheInuse)).
		add(adapter.NewUpdateRequestMetricGauge("MCacheSys", stats.Runtime.MCacheSys)).
		add(adapter.NewUpdateRequestMetricGauge("MSpanInuse", stats.Runtime.MSpanInuse)).
		add(adapter.NewUpdateRequestMetricGauge("MSpanSys", stats.Runtime.MSpanSys)).
		add(adapter.NewUpdateRequestMetricGauge("Mallocs", stats.Runtime.Mallocs)).
		add(adapter.NewUpdateRequestMetricGauge("NextGC", stats.Runtime.NextGC)).
		add(adapter.NewUpdateRequestMetricGauge("NumForcedGC", stats.Runtime.NumForcedGC)).
		add(adapter.NewUpdateRequestMetricGauge("NumGC", stats.Runtime.NumGC)).
		add(adapter.NewUpdateRequestMetricGauge("OtherSys", stats.Runtime.OtherSys)).
		add(adapter.NewUpdateRequestMetricGauge("PauseTotalNs", stats.Runtime.PauseTotalNs)).
		add(adapter.NewUpdateRequestMetricGauge("StackInuse", stats.Runtime.StackInuse)).
		add(adapter.NewUpdateRequestMetricGauge("StackSys", stats.Runtime.StackSys)).
		add(adapter.NewUpdateRequestMetricGauge("Sys", stats.Runtime.Sys)).
		add(adapter.NewUpdateRequestMetricGauge("TotalAlloc", stats.Runtime.TotalAlloc))
	// отправляем метрики пакета gopsutil
	sender.
		add(adapter.NewUpdateRequestMetricGauge("TotalMemory", stats.PS.TotalMemory)).
		add(adapter.NewUpdateRequestMetricGauge("FreeMemory", stats.PS.FreeMemory)).
		add(adapter.NewUpdateRequestMetricGauge("CPUutilization1", stats.PS.CPUutilization1))
	// отправляем обновляемое произвольное значение
	sender.
		add(adapter.NewUpdateRequestMetricGauge("RandomValue", stats.RandomValue))
	// отправляем счётчик обновления метрик пакета runtime
	sender.
		add(adapter.NewUpdateRequestMetricCounter("PollCount", stats.PollCount))

	close(sender.batch)

	return sender.err
}

func (hs *httpSender) worker(ctx context.Context) {
	data := make([]adapter.RequestMetric, 0, len(hs.batch))

	for r := range hs.batch {
		data = append(data, r)
	}

	hs.err = hs.doSend(ctx, data)
}
