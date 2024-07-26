package metric

import (
	"strconv"
)

type (
	// Counter тип метрики.
	Counter int64 
)

// Kind возвращает строку, тип метрики.
func (c Counter) Kind() string {
	return string(KindCounter)
}

// String возвращает строковое представление значения.
func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

// IsCounter возвращает true, если метрика имеет тип counter.
func (c Counter) IsCounter() bool {
	return true
}

// IsGauge возвращает true, если метрика имеет тип gauge.
func (c Counter) IsGauge() bool {
	return false
}

// ToCounter преобразует строковое представление значения метрики в значение с типом counter.
// 
// Параметры:
//	- value - строковое представление значения метрики.
//
// Возвращаемое значение:
//	- Counter - значение с типом counter.
//	- error - ошибка, если преобразование не удалось.
//
func ToCounter(value string) (Counter, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return Counter(val), nil
}
