package metricservice

import (
	"fmt"
	"strconv"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

type metricService struct {
	stor storage.Storage
}

func New(stor storage.Storage) metricService {
	return metricService{stor}
}

func (s metricService) Save(name, kind, value string) error {
	metricKind, err := metric.GetKind(kind)
	if err != nil {
		return err
	}

	record := storage.NewRecord(name)

	switch metricKind {
	case metric.KindGauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		record.SetValue(metric.Gauge(val))
	case metric.KindCounter:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		record.SetValue(metric.Counter(val))
	default:
		return metric.ErrorInvalidMetricKind
	}

	fmt.Println("name:", name, record)

	return s.stor.Push(name, record)
}
