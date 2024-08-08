// Package adapter методы для конвертации внутреней модели metric в модель используему в запросах.
package adapter

import "github.com/a-x-a/go-metric/internal/models/metric"

func NewUpdateRequestMetricCounter(name string, value metric.Counter) metric.RequestMetric {
	val := int64(value)

	return metric.RequestMetric{
		ID:    name,
		MType: value.Kind(),
		Delta: &val,
	}
}

func NewUpdateRequestMetricGauge(name string, value metric.Gauge) metric.RequestMetric {
	val := float64(value)

	return metric.RequestMetric{
		ID:    name,
		MType: value.Kind(),
		Value: &val,
	}
}

func NewGetRequestMetricCounter(name string) metric.RequestMetric {
	return metric.RequestMetric{
		ID:    name,
		MType: string(metric.KindCounter),
	}
}

func NewGetRequestMetricGauge(name string) metric.RequestMetric {
	return metric.RequestMetric{
		ID:    name,
		MType: string(metric.KindGauge),
	}
}
