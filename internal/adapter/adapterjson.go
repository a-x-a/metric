package adapter

import "github.com/a-x-a/go-metric/internal/models/metric"

type RequestMetric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewUpdateRequestMetricCounter(name string, value metric.Counter) RequestMetric {
	val := int64(value)

	return RequestMetric{
		ID:    name,
		MType: value.Kind(),
		Delta: &val,
	}
}

func NewUpdateRequestMetricGauge(name string, value metric.Gauge) RequestMetric {
	val := float64(value)

	return RequestMetric{
		ID:    name,
		MType: value.Kind(),
		Value: &val,
	}
}

func NewGetRequestMetricCounter(name string) RequestMetric {
	return RequestMetric{
		ID:    name,
		MType: string(metric.KindCounter),
	}
}

func NewGetRequestMetricGauge(name string) RequestMetric {
	return RequestMetric{
		ID:    name,
		MType: string(metric.KindGauge),
	}
}
