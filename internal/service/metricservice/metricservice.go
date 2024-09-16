package metricservice

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	metricService struct {
		storage storage.Storage
		logger  *zap.Logger
	}

	dbStorage interface {
		Ping(ctx context.Context) error
	}
)

var ErrNotSupportedMethod = errors.New("storage doesn't support method")

func New(stor storage.Storage, logger *zap.Logger) *metricService {
	return &metricService{
		storage: stor,
		logger:  logger,
	}
}

func (s *metricService) Push(name, kind, value string) error {
	metricKind, err := metric.GetKind(kind)
	if err != nil {
		return err
	}

	record, err := storage.NewRecord(name)
	if err != nil {
		return err
	}

	switch metricKind {
	case metric.KindGauge:
		val, err := metric.ToGauge(value)
		if err != nil {
			return err
		}
		record.SetValue(val)
	case metric.KindCounter:
		val, err := metric.ToCounter(value)
		if err != nil {
			return err
		}
		if v, ok := s.storage.Get(name); ok {
			if oldVal, ok := v.GetValue().(metric.Counter); ok {
				val += oldVal
			}
		}
		record.SetValue(val)
	default:
		return metric.ErrorInvalidMetricKind
	}

	return s.storage.Push(name, record)
}

func (s *metricService) PushCounter(name string, value metric.Counter) (metric.Counter, error) {
	record, err := storage.NewRecord(name)
	if err != nil {
		return 0, err
	}

	if v, ok := s.storage.Get(name); ok {
		if oldVal, ok := v.GetValue().(metric.Counter); ok {
			value += oldVal
		}
	}
	record.SetValue(value)

	return value, s.storage.Push(name, record)
}

func (s *metricService) PushGauge(name string, value metric.Gauge) (metric.Gauge, error) {
	record, err := storage.NewRecord(name)
	if err != nil {
		return 0, err
	}

	record.SetValue(value)

	err = s.storage.Push(name, record)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s metricService) Get(name, kind string) (string, error) {
	if _, err := metric.GetKind(kind); err != nil {
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

func (s metricService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbStorage, ok := s.storage.(dbStorage)

	if !ok {
		return ErrNotSupportedMethod
	}

	return dbStorage.Ping(ctx)
}
