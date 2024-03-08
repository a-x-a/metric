package metricservice

import (
	"strconv"

	"github.com/a-x-a/go-metric/internal/model/metric"
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

	rec := storage.NewRecord(name)

	switch metricKind {
	case metric.KindGauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		rec.SetValue(metric.Gauge(val))
	case metric.KindCounter:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		rec.SetValue(metric.Counter(val))
	default:
		return metric.ErrorInvalidMetricKind
	}

	return s.stor.Save(name, rec)
}
