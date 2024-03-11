package metricservice

import (
	"fmt"
	"strconv"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	MetricService interface {
		Push(name, kind, value string) error
		Get(name, kind string) (string, error)
		GetAll() []storage.Record
	}

	metricService struct {
		storage storage.Storage
	}
)

func New(stor storage.Storage) metricService {
	return metricService{stor}
}

func (s metricService) Push(name, kind, value string) error {
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

	return s.storage.Push(name, record)
}

func (s metricService) Get(name, kind string) (string, error) {
	_, err := metric.GetKind(kind)
	if err != nil {
		return "", err
	}

	record, ok := s.storage.Get(name)
	if !ok {
		return "", metric.ErrorMetricNotFound
	}

	value := record.GetValue().String()
	return value, nil
}

func (s metricService) GetAll() []storage.Record {
	records := s.storage.GetAll()

	return records
}
