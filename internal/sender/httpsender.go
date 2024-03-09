package sender

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

type httpSender struct {
	baseURL string
	client  *http.Client
	err     error
}

func NewSender(serverAddress string) httpSender {
	baseURL := fmt.Sprintf("http://%s", serverAddress)
	client := &http.Client{Timeout: 2 * time.Second}

	return httpSender{baseURL: baseURL, client: client, err: nil}
}

func (hs *httpSender) doSend(req string) *httpSender {
	resp, err := hs.client.Post(req, "Content-Type: text/plain", nil)
	if err != nil {
		hs.err = err
		return hs
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		hs.err = err
		return hs
	}

	if resp.StatusCode != http.StatusOK {
		hs.err = fmt.Errorf("metrics send failed: (%d)", resp.StatusCode)
		return hs
	}

	return hs
}

func (hs *httpSender) exportGauge(name string, value metric.Gauge) *httpSender {
	if hs.err != nil {
		return hs
	}

	req := fmt.Sprintf("%s/update/gauge/%s/%f", hs.baseURL, name, value)

	return hs.doSend(req)
}

func (hs *httpSender) exportCounter(name string, value metric.Counter) *httpSender {
	if hs.err != nil {
		return hs
	}

	req := fmt.Sprintf("%s/update/counter/%s/%d", hs.baseURL, name, value)

	return hs.doSend(req)
}

func SendMetrics(serverAddress string, stats metric.Metrics) error {
	sender := NewSender(serverAddress)

	// отправляем метрики пакета runtime
	sender.
		exportGauge("Alloc", stats.Memory.Alloc).
		exportGauge("BuckHashSys", stats.Memory.BuckHashSys).
		exportGauge("Frees", stats.Memory.Frees).
		exportGauge("GCCPUFraction", stats.Memory.GCCPUFraction).
		exportGauge("GCSys", stats.Memory.GCSys).
		exportGauge("HeapAlloc", stats.Memory.HeapAlloc).
		exportGauge("HeapIdle", stats.Memory.HeapIdle).
		exportGauge("HeapInuse", stats.Memory.HeapInuse).
		exportGauge("HeapObjects", stats.Memory.HeapObjects).
		exportGauge("HeapReleased", stats.Memory.HeapReleased).
		exportGauge("HeapSys", stats.Memory.HeapSys).
		exportGauge("LastGC", stats.Memory.LastGC).
		exportGauge("Lookups", stats.Memory.Lookups).
		exportGauge("MCacheInuse", stats.Memory.MCacheInuse).
		exportGauge("MCacheSys", stats.Memory.MCacheSys).
		exportGauge("MSpanInuse", stats.Memory.MSpanInuse).
		exportGauge("MSpanSys", stats.Memory.MSpanSys).
		exportGauge("Mallocs", stats.Memory.Mallocs).
		exportGauge("NextGC", stats.Memory.NextGC).
		exportGauge("NumForcedGC", stats.Memory.NumForcedGC).
		exportGauge("NumGC", stats.Memory.NumGC).
		exportGauge("OtherSys", stats.Memory.OtherSys).
		exportGauge("PauseTotalNs", stats.Memory.PauseTotalNs).
		exportGauge("StackInuse", stats.Memory.StackInuse).
		exportGauge("StackSys", stats.Memory.StackSys).
		exportGauge("Sys", stats.Memory.Sys).
		exportGauge("TotalAlloc", stats.Memory.TotalAlloc)

	// отправляем обновляемое произвольное значение
	sender.exportGauge("RandomValue", stats.RandomValue)
	// отправляем счётчик обновления метрик пакета runtime
	sender.exportCounter("PollCount", stats.PollCount)

	return sender.err
}
