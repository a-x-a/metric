package metric

import (
	"errors"
)

type (
	MetricKind string
)

const (
	// типы метрик
	KindGauge, KindCounter MetricKind = "gauge", "counter"
)

var (
	// metricTypes - строковое представление допустимых типов метрик
	metricKinds = map[string]MetricKind{"gauge": KindGauge, "counter": KindCounter}
	// ErrorInvalidMetricType - не корректный тип метрики.
	ErrorInvalidMetricKind = errors.New("model: не корректный тип метрики")
)

// GetKind - возвращает корректный тип метрики для строкового представления
// Если передан не корректный тип метрики, то возвращает ошибку ErrorInvalidMetricType
func GetKind(kindRaw string) (MetricKind, error) {
	if v, ok := metricKinds[kindRaw]; ok {
		return v, nil
	}
	return MetricKind(""), ErrorInvalidMetricKind
}
