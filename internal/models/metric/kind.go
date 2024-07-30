package metric

import (
	"errors"
)

type (
	// MetricKind тип метрики.
	MetricKind string
)

const (
	// KindGauge, KindCounter типы метрик.
	KindGauge, KindCounter MetricKind = "gauge", "counter"
)

var (
	// metricTypes строковое представление допустимых типов метрик.
	metricKinds = map[string]MetricKind{"gauge": KindGauge, "counter": KindCounter}
	// ErrorInvalidMetricType ошибка, если передан не корректный тип метрики.
	ErrorInvalidMetricKind = errors.New("model: не корректный тип метрики")
)

// GetKind возвращает корректный тип метрики для строкового представления.
//
// Параметры:
//   - kindRaw - строковое представление типа метрики.
//
// Возвращаемое значение:
//   - MetricKind - корректный тип метрики.
//   - error - ошибка, если передан не корректный тип метрики.
func GetKind(kindRaw string) (MetricKind, error) {
	if v, ok := metricKinds[kindRaw]; ok {
		return v, nil
	}

	return MetricKind(""), ErrorInvalidMetricKind
}
