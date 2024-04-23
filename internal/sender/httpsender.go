package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-x-a/go-metric/internal/adapter"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

type httpSender struct {
	baseURL string
	client  *http.Client
	batch   []adapter.RequestMetric
	err     error
}

func NewSender(serverAddress string, timeout time.Duration) httpSender {
	baseURL := fmt.Sprintf("http://%s", serverAddress)
	client := &http.Client{Timeout: timeout}

	return httpSender{
		baseURL: baseURL,
		client:  client,
		batch:   make([]adapter.RequestMetric, 0),
		err:     nil,
	}
}

func (hs *httpSender) doSend(ctx context.Context) error {
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

	hs.err = hs.doSend(ctx)

	return hs
}

func (hs *httpSender) Add(rm adapter.RequestMetric) *httpSender {
	if hs.err != nil {
		return hs
	}

	hs.batch = append(hs.batch, rm)

	return hs
}

func SendMetrics(ctx context.Context, serverAddress string, timeout time.Duration, stats metric.Metrics) error {
	sender := NewSender(serverAddress, timeout)

	// отправляем метрики пакета runtime
	sender.
		Add(adapter.NewUpdateRequestMetricGauge("Alloc", stats.Memory.Alloc)).
		Add(adapter.NewUpdateRequestMetricGauge("BuckHashSys", stats.Memory.BuckHashSys)).
		Add(adapter.NewUpdateRequestMetricGauge("Frees", stats.Memory.Frees)).
		Add(adapter.NewUpdateRequestMetricGauge("GCCPUFraction", stats.Memory.GCCPUFraction)).
		Add(adapter.NewUpdateRequestMetricGauge("GCSys", stats.Memory.GCSys)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapAlloc", stats.Memory.HeapAlloc)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapIdle", stats.Memory.HeapIdle)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapInuse", stats.Memory.HeapInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapObjects", stats.Memory.HeapObjects)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapReleased", stats.Memory.HeapReleased)).
		Add(adapter.NewUpdateRequestMetricGauge("HeapSys", stats.Memory.HeapSys)).
		Add(adapter.NewUpdateRequestMetricGauge("LastGC", stats.Memory.LastGC)).
		Add(adapter.NewUpdateRequestMetricGauge("Lookups", stats.Memory.Lookups)).
		Add(adapter.NewUpdateRequestMetricGauge("MCacheInuse", stats.Memory.MCacheInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("MCacheSys", stats.Memory.MCacheSys)).
		Add(adapter.NewUpdateRequestMetricGauge("MSpanInuse", stats.Memory.MSpanInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("MSpanSys", stats.Memory.MSpanSys)).
		Add(adapter.NewUpdateRequestMetricGauge("Mallocs", stats.Memory.Mallocs)).
		Add(adapter.NewUpdateRequestMetricGauge("NextGC", stats.Memory.NextGC)).
		Add(adapter.NewUpdateRequestMetricGauge("NumForcedGC", stats.Memory.NumForcedGC)).
		Add(adapter.NewUpdateRequestMetricGauge("NumGC", stats.Memory.NumGC)).
		Add(adapter.NewUpdateRequestMetricGauge("OtherSys", stats.Memory.OtherSys)).
		Add(adapter.NewUpdateRequestMetricGauge("PauseTotalNs", stats.Memory.PauseTotalNs)).
		Add(adapter.NewUpdateRequestMetricGauge("StackInuse", stats.Memory.StackInuse)).
		Add(adapter.NewUpdateRequestMetricGauge("StackSys", stats.Memory.StackSys)).
		Add(adapter.NewUpdateRequestMetricGauge("Sys", stats.Memory.Sys)).
		Add(adapter.NewUpdateRequestMetricGauge("TotalAlloc", stats.Memory.TotalAlloc))

	// отправляем обновляемое произвольное значение
	sender.
		Add(adapter.NewUpdateRequestMetricGauge("RandomValue", stats.RandomValue))
	// отправляем счётчик обновления метрик пакета runtime
	sender.
		Add(adapter.NewUpdateRequestMetricCounter("PollCount", stats.PollCount))

	return sender.Send(ctx).err
}
