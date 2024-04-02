package metric

import (
	"strconv"
)

type (
	Gauge float64
)

func (g Gauge) Kind() string {
	return string(KindGauge)
}

func (g Gauge) String() string {
	if g == 0 {
		return strconv.FormatFloat(float64(g), 'f', 3, 64)
	}
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func (g Gauge) IsCounter() bool {
	return false
}

func (g Gauge) IsGauge() bool {
	return true
}

func ToGauge(value string) (Gauge, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return Gauge(val), nil
}
