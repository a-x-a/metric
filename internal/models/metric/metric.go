package metric

import (
	"context"
	"errors"
	"math/rand"

	"golang.org/x/sync/errgroup"
)

type (
	Metric interface {
		Kind() string   // Kind - возвращает тип метрики.
		String() string // Stringer.
		IsCounter() bool
		IsGauge() bool
	}

	Metrics struct {
		// метрики пакета runtime.
		Runtime RuntimeMetrics
		// метрики пакета gopsutil.
		PS PSMetrics
		// дополнительные метрики.
		// PollCount - счётчик, увеличивающийся на 1 при каждом обновлении метрики из пакета runtime.
		PollCount Counter
		// RandomValue - обновляемое произвольное значение.
		RandomValue Gauge
	}
)

var (
	// ErrorMetricNameIsNull - не указано имя метрики.
	ErrorMetricNameIsNull = errors.New("metrics: ошибка cоздания метрики, не указано име метрики")
	// ErrorMetricNotFound - метрика не найдена.
	ErrorMetricNotFound = errors.New("metrics: метрика не найдена")
)

func (m *Metrics) Poll(ctx context.Context) error {
	m.PollCount += 1
	m.RandomValue = Gauge(rand.Float64())

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		m.Runtime.Poll()

		return nil
	})

	g.Go(func() error {
		return m.PS.Poll()
	})

	return g.Wait()
}
