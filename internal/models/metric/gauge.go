package metric

import (
	"strconv"
)

type (
	// Gauge тип метрики.
	Gauge float64
)

// Kind возвращает строку, тип метрики.
func (g Gauge) Kind() string {
	return string(KindGauge)
}

// String возвращает строковое представление значения.
func (g Gauge) String() string {
	if g == 0 {
		return strconv.FormatFloat(float64(g), 'f', 3, 64)
	}
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

// IsCounter возвращает true, если метрика имеет тип counter.
func (g Gauge) IsCounter() bool {
	return false
}

// IsGauge возвращает true, если метрика имеет тип gauge.
func (g Gauge) IsGauge() bool {
	return true
}

// ToGauge преобразует строковое представление значения метрики в значение с типом gauge.
//
// Параметры:
//   - value - строковое представление значения метрики.
//
// Возвращаемое значение:
//   - Gauge - значение с типом gauge.
//   - error - ошибка, если преобразование не удалось.
func ToGauge(value string) (Gauge, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return Gauge(val), nil
}
