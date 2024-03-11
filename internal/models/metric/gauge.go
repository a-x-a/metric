package metric

import "fmt"

type (
	// Gauge interface {
	// 	Get() float64
	// 	Set(v float64) error
	// }

	Gauge float64
)

func (g Gauge) Kind() string {
	return string(KindGauge)
}

func (g Gauge) String() string {
	return fmt.Sprintf("%.3f", g)
}
