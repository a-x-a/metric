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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
		if v, err := s.storage.Get(ctx, name); err == nil {
			if oldVal, ok := v.GetValue().(metric.Counter); ok {
				val += oldVal
			}
		}
		record.SetValue(val)
	}

	return s.storage.Push(ctx, name, record)
}

func (s *metricService) PushCounter(name string, value metric.Counter) (metric.Counter, error) {
	record, err := storage.NewRecord(name)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if v, err := s.storage.Get(ctx, name); err == nil {
		if oldVal, ok := v.GetValue().(metric.Counter); ok {
			value += oldVal
		}
	}
	record.SetValue(value)

	return value, s.storage.Push(ctx, name, record)
}

func (s *metricService) PushGauge(name string, value metric.Gauge) (metric.Gauge, error) {
	record, err := storage.NewRecord(name)
	if err != nil {
		return 0, err
	}

	record.SetValue(value)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = s.storage.Push(ctx, name, record)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s metricService) PushBatch(ctx context.Context, records []storage.Record) error {
	data := make([]storage.Record, 0, len(records))
	cache := make(map[string]int)
	counters := make(map[string]metric.Counter)

	for _, v := range records {
		id := v.GetName()
		value := v.GetValue()

		if i, ok := cache[id]; ok {
			if value.IsGauge() {
				data[i].SetValue(value)
				continue
			}

			if value.IsCounter() {
				if oldValue, ok := counters[id]; ok {
					value = oldValue + value.(metric.Counter)
				}
				counters[id] = value.(metric.Counter)
				data[i].SetValue(value)
				continue
			}
		}

		record, err := storage.NewRecord(id)
		if err != nil {
			return err
		}

		if value.IsCounter() {
			storRecord, err := s.Get(id, value.Kind())
			if err != nil && !errors.Is(err, metric.ErrorMetricNotFound) {
				return err
			}

			if storRecord != nil {
				oldValue := storRecord.GetValue()
				value = oldValue.(metric.Counter) + value.(metric.Counter)
			}
			counters[id] = value.(metric.Counter)
		}

		record.SetValue(value)
		cache[id] = len(data)
		data = append(data, record)
	}

	return s.storage.PushBatch(ctx, data)
}

func (s metricService) Get(name, kind string) (*storage.Record, error) {
	if _, err := metric.GetKind(kind); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	record, err := s.storage.Get(ctx, name)
	if err != nil {
		return nil, metric.ErrorMetricNotFound
	}

	return record, nil
}

func (s metricService) GetAll() []storage.Record {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	records, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil
	}

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
