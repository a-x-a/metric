package metric

import (
	"context"
	"errors"
	"math/rand"

	"golang.org/x/sync/errgroup"
)

type (
	// Metric методы метрик.
	Metric interface {
		// Kind возвращает тип метрики.
		Kind() string
		// String возвращает строковое представление значения метрики.
		String() string
		// IsCounter возвращает true если метрика является счётчиком.
		IsCounter() bool
		// IsGauge возвращает true если метрика является датчиком.
		IsGauge() bool
	}

	// Metrics структура метрик.
	Metrics struct {
		// Runtime метрики пакета runtime.
		Runtime RuntimeMetrics
		// PS метрики пакета gopsutil.
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

// Poll обновление метрик.
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
