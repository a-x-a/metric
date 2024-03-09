package metric

import (
	"errors"
)

type (
	// MetricKind string

	Metric interface {
		Kind() string // Kind - возвращает тип метрики
	}
)

// const (
// 	// типы метрик
// 	KindGauge, KindCounter MetricKind = "gauge", "counter"
// )

var (
	// ErrorMetricNameIsNull - не указано имя метрики
	ErrorMetricNameIsNull = errors.New("model: ошибка cоздания метрики, не указано име метрики")
)
