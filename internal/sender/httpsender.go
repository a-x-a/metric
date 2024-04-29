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
	"github.com/a-x-a/go-metric/internal/signer"
)

type httpSender struct {
	baseURL string
	client  *http.Client
	signer  *signer.Signer
	batch   chan adapter.RequestMetric
	err     error
}

func NewHTTPSender(serverAddress string, timeout time.Duration, key string) httpSender {
	baseURL := fmt.Sprintf("http://%s", serverAddress)
	client := &http.Client{Timeout: timeout}
	sgnr := signer.New(key)

	return httpSender{
		baseURL: baseURL,
		client:  client,
		signer:  sgnr,
		batch:   make(chan adapter.RequestMetric, 1024),
		err:     nil,
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

	var buf bytes.Buffer

	zw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return err
	}

	if _, err := zw.Write(data); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hs.baseURL+"/updates", &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	if hs.signer != nil {
		hash, err := hs.signer.Hash(data)
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

func (hs *httpSender) do(ctx context.Context) error {
	data, err := json.Marshal(hs.batch)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	zw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return err
	}

	if _, err := zw.Write(data); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hs.baseURL+"/updates", &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	if hs.signer != nil {
		hash, err := hs.signer.Hash(data)
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

func (hs *httpSender) Send(ctx context.Context) *httpSender {
	if hs.err != nil {
		return hs
	}

	if len(hs.batch) == 0 {
		hs.err = fmt.Errorf("metrics send: batch is empty")
		return hs
	}

	hs.err = hs.do(ctx)

	return hs
}

func (hs *httpSender) Add(rm adapter.RequestMetric) *httpSender {
	if hs.err != nil {
		return hs
	}

	hs.batch <- rm

	return hs
}

func SendMetrics(ctx context.Context, serverAddress string, timeout time.Duration, key string, rateLimit int, stats metric.Metrics) error {
	sender := NewHTTPSender(serverAddress, timeout, key)

	for i := 0; i < rateLimit; i++ {
		go sender.worker(ctx)
	}

	// отправляем метрики пакета runtime
	sender.
		Add(adapter.NewUpdateRequestMetricGauge("Alloc", stats.Runtime.Alloc)).
		Add(adapter.NewUpdateRequestMetricGauge("BuckHashSys", stats.Runtime.BuckHashSys)).
		Add(adapter.NewUpdateRequestMetricGauge("Frees", stats.Runtime.Frees)).
		Add(adapter.NewUpdateRequestMetricGauge("GCCPUFraction", stats.Runtime.GCCPUFraction)).
		Add(adapter.NewUpdateRequestMetricGauge("GCSys", stats.Runtime.GCSys)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapAlloc", stats.Runtime.HeapAlloc)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapIdle", stats.Runtime.HeapIdle)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapInuse", stats.Runtime.HeapInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapObjects", stats.Runtime.HeapObjects)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapReleased", stats.Runtime.HeapReleased)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapSys", stats.Runtime.HeapSys)).
		Add(adapter.NewUpdateRequestMetricGauge("LastGC", stats.Runtime.LastGC)).
		Add(adapter.NewUpdateRequestMetricGauge("Lookups", stats.Runtime.Lookups)).
		Add(adapter.NewUpdateRequestMetricGauge("MCacheInuse", stats.Runtime.MCacheInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("MCacheSys", stats.Runtime.MCacheSys)).
		Add(adapter.NewUpdateRequestMetricGauge("MSpanInuse", stats.Runtime.MSpanInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("MSpanSys", stats.Runtime.MSpanSys)).
		Add(adapter.NewUpdateRequestMetricGauge("Mallocs", stats.Runtime.Mallocs)).
		Add(adapter.NewUpdateRequestMetricGauge("NextGC", stats.Runtime.NextGC)).
		Add(adapter.NewUpdateRequestMetricGauge("NumForcedGC", stats.Runtime.NumForcedGC)).
		Add(adapter.NewUpdateRequestMetricGauge("NumGC", stats.Runtime.NumGC)).
		Add(adapter.NewUpdateRequestMetricGauge("OtherSys", stats.Runtime.OtherSys)).
		Add(adapter.NewUpdateRequestMetricGauge("PauseTotalNs", stats.Runtime.PauseTotalNs)).
		Add(adapter.NewUpdateRequestMetricGauge("StackInuse", stats.Runtime.StackInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("StackSys", stats.Runtime.StackSys)).
		Add(adapter.NewUpdateRequestMetricGauge("Sys", stats.Runtime.Sys)).
		Add(adapter.NewUpdateRequestMetricGauge("TotalAlloc", stats.Runtime.TotalAlloc))
	// отправляем метрики пакета gopsutil
	sender.
		Add(adapter.NewUpdateRequestMetricGauge("TotalMemory", stats.PS.TotalMemory)).
		Add(adapter.NewUpdateRequestMetricGauge("FreeMemory", stats.PS.FreeMemory)).
		Add(adapter.NewUpdateRequestMetricGauge("CPUutilization1", stats.PS.CPUutilization1))
	// отправляем обновляемое произвольное значение
	sender.
		Add(adapter.NewUpdateRequestMetricGauge("RandomValue", stats.RandomValue))
	// отправляем счётчик обновления метрик пакета runtime
	sender.
		Add(adapter.NewUpdateRequestMetricCounter("PollCount", stats.PollCount))

	close(sender.batch)

	return sender.err
}

func (hs *httpSender) worker(ctx context.Context) {
	data := make([]adapter.RequestMetric, 0)
	for {
		select {
		case r, ok := <-hs.batch:
			if ok {
				data = append(data, r)
				continue
			}

			if len(data) != 0 {
				hs.err = hs.doSend(ctx, data)
			}

			return
		}
	}
}
