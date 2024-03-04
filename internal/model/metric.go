package model

import (
	"errors"
)

type (
	MetricGuage   float64
	MetricCounter int64
	MetricType    int

	metricValue struct {
		guage   MetricGuage
		counter MetricCounter
	}

	metric struct {
		name       string
		metricType MetricType
		value      metricValue
	}
)

const (
	Guage MetricType = iota + 1
	Counter
)

var (
	// metricTypes - строковое представление допустимых типов метрик.
	metricTypes = [...]string{"", "guage", "counter"}

	// ErroMetricNameIsNull - не указано имя метрики
	ErroMetricNameIsNull = errors.New("model: ошибка cоздания метрики, не указано име метрики")
	// ErroInvalidMetricType - не корректный тип метрики.
	ErroInvalidMetricType = errors.New("model: ошибка cоздания метрики, не корректный тип метрики")
)

// String - возвращает строковое представление типа метрики.
func (mt MetricType) String() string {
	if mt.isValid() {
		return metricTypes[mt]
	}
	return ""
}

// isValid - проверяет правильность типа метрики
func (mt MetricType) isValid() bool {
	return int(mt) > 0 && int(mt) < len(metricTypes)
}

// NewMetric - возвращает указатель на новый объект метрики.
// Проверяет правильность указания имени и типа метрики.
// Вслучае получения не корректных значений для создания метрики, возвращает ошибку.
func NewMetric(name string, metricType int) (*metric, error) {
	if name == "" {
		return nil, ErroMetricNameIsNull
	}
	if !MetricType(metricType).isValid() {
		return nil, ErroInvalidMetricType
	}

	return &metric{
		name:       name,
		metricType: MetricType(metricType),
		value: metricValue{
			guage:   0,
			counter: 0,
		},
	}, nil
}
