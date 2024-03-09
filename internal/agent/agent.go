package agent

import (
	"runtime"

	"github.com/a-x-a/go-metric/internal/model/metric"
)

type RuntimeMetrics struct {
	Alloc         metric.Gauge
	BuckHashSys   metric.Gauge
	Frees         metric.Gauge
	GCCPUFraction metric.Gauge
	GCSys         metric.Gauge
	HeapAlloc     metric.Gauge
	HeapIdle      metric.Gauge
	HeapInuse     metric.Gauge
	HeapObjects   metric.Gauge
	HeapReleased  metric.Gauge
	HeapSys       metric.Gauge
	LastGC        metric.Gauge
	Lookups       metric.Gauge
	MCacheInuse   metric.Gauge
	MCacheSys     metric.Gauge
	MSpanInuse    metric.Gauge
	MSpanSys      metric.Gauge
	Mallocs       metric.Gauge
	NextGC        metric.Gauge
	NumForcedGC   metric.Gauge
	NumGC         metric.Gauge
	OtherSys      metric.Gauge
	PauseTotalNs  metric.Gauge
	StackInuse    metric.Gauge
	StackSys      metric.Gauge
	Sys           metric.Gauge
	TotalAlloc    metric.Gauge
}

// Poll - обновляет значения показателей метрик
func (m *RuntimeMetrics) Poll() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	m.Alloc = metric.Gauge(stats.Alloc)
	m.BuckHashSys = metric.Gauge(stats.BuckHashSys)
	m.Frees = metric.Gauge(stats.Frees)
	m.GCCPUFraction = metric.Gauge(stats.GCCPUFraction)
	m.GCSys = metric.Gauge(stats.GCSys)
	m.HeapAlloc = metric.Gauge(stats.HeapAlloc)
	m.HeapIdle = metric.Gauge(stats.HeapIdle)
	m.HeapInuse = metric.Gauge(stats.HeapInuse)
	m.HeapObjects = metric.Gauge(stats.HeapObjects)
	m.HeapReleased = metric.Gauge(stats.HeapReleased)
	m.HeapSys = metric.Gauge(stats.HeapSys)
	m.LastGC = metric.Gauge(stats.LastGC)
	m.Lookups = metric.Gauge(stats.Lookups)
	m.MCacheInuse = metric.Gauge(stats.MCacheInuse)
	m.MCacheSys = metric.Gauge(stats.MCacheSys)
	m.MSpanInuse = metric.Gauge(stats.MSpanInuse)
	m.MSpanSys = metric.Gauge(stats.MSpanSys)
	m.Mallocs = metric.Gauge(stats.Mallocs)
	m.NextGC = metric.Gauge(stats.NextGC)
	m.NumForcedGC = metric.Gauge(stats.NumForcedGC)
	m.NumGC = metric.Gauge(stats.NumGC)
	m.OtherSys = metric.Gauge(stats.OtherSys)
	m.PauseTotalNs = metric.Gauge(stats.PauseTotalNs)
	m.StackInuse = metric.Gauge(stats.StackInuse)
	m.StackSys = metric.Gauge(stats.StackSys)
	m.Sys = metric.Gauge(stats.Sys)
	m.TotalAlloc = metric.Gauge(stats.TotalAlloc)
}

func SendMetric() error {

	return nil
}
