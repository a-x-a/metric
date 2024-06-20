package metric

import (
	"runtime"
)

type (
	//RuntimeMetrics метрики пакета runtime.
	RuntimeMetrics struct {
		Alloc         Gauge
		BuckHashSys   Gauge
		Frees         Gauge
		GCCPUFraction Gauge
		GCSys         Gauge
		HeapAlloc     Gauge
		HeapIdle      Gauge
		HeapInuse     Gauge
		HeapObjects   Gauge
		HeapReleased  Gauge
		HeapSys       Gauge
		LastGC        Gauge
		Lookups       Gauge
		MCacheInuse   Gauge
		MCacheSys     Gauge
		MSpanInuse    Gauge
		MSpanSys      Gauge
		Mallocs       Gauge
		NextGC        Gauge
		NumForcedGC   Gauge
		NumGC         Gauge
		OtherSys      Gauge
		PauseTotalNs  Gauge
		StackInuse    Gauge
		StackSys      Gauge
		Sys           Gauge
		TotalAlloc    Gauge
	}
)

// Poll обновляет значения показателей метрик.
func (rm *RuntimeMetrics) Poll() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	rm.Alloc = Gauge(stats.Alloc)
	rm.BuckHashSys = Gauge(stats.BuckHashSys)
	rm.Frees = Gauge(stats.Frees)
	rm.GCCPUFraction = Gauge(stats.GCCPUFraction)
	rm.GCSys = Gauge(stats.GCSys)
	rm.HeapAlloc = Gauge(stats.HeapAlloc)
	rm.HeapIdle = Gauge(stats.HeapIdle)
	rm.HeapInuse = Gauge(stats.HeapInuse)
	rm.HeapObjects = Gauge(stats.HeapObjects)
	rm.HeapReleased = Gauge(stats.HeapReleased)
	rm.HeapSys = Gauge(stats.HeapSys)
	rm.LastGC = Gauge(stats.LastGC)
	rm.Lookups = Gauge(stats.Lookups)
	rm.MCacheInuse = Gauge(stats.MCacheInuse)
	rm.MCacheSys = Gauge(stats.MCacheSys)
	rm.MSpanInuse = Gauge(stats.MSpanInuse)
	rm.MSpanSys = Gauge(stats.MSpanSys)
	rm.Mallocs = Gauge(stats.Mallocs)
	rm.NextGC = Gauge(stats.NextGC)
	rm.NumForcedGC = Gauge(stats.NumForcedGC)
	rm.NumGC = Gauge(stats.NumGC)
	rm.OtherSys = Gauge(stats.OtherSys)
	rm.PauseTotalNs = Gauge(stats.PauseTotalNs)
	rm.StackInuse = Gauge(stats.StackInuse)
	rm.StackSys = Gauge(stats.StackSys)
	rm.Sys = Gauge(stats.Sys)
	rm.TotalAlloc = Gauge(stats.TotalAlloc)
}
