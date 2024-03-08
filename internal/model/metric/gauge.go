package metric

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
