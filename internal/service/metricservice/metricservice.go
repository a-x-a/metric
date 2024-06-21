package metricservice

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

type (
	// MetricService сервис сбора метрик.
	MetricService struct {
		storage stor
		logger  *zap.Logger
	}
	
	// StorageWithPing хранилище, поддерживающее метод Ping.
	StorageWithPing interface {
		Ping(ctx context.Context) error
	}

	// stor основные методы хранилища
	stor interface {
		Push(ctx context.Context, name string, record storage.Record) error
		PushBatch(ctx context.Context, records []storage.Record) error
		Get(ctx context.Context, name string) (*storage.Record, error)
		GetAll(ctx context.Context) ([]storage.Record, error)
	}
)

// ErrNotSupportedMethod ошибка, не поддерживаемый метод
var ErrNotSupportedMethod = errors.New("storage doesn't support method")

// New создает новый экземпляр сервиса сбора метрик.
//
// Параметры:
// - stor - хранилище метрик;
// - logger - логгер.
//
// Возвращаемое значение:
// - *MetricService - сервис сбора метрик.
//
func New(stor storage.Storage, logger *zap.Logger) *MetricService {
	return &MetricService{
		storage: stor,
		logger:  logger,
	}
}

// Push добавляет метрику с указанным именем, типом и значением.
//
// Параметры:
// - ctx - контекст;
// - name - имя метрики;
// - kind - тип метрики;
// - value - значение метрики.
//
// Возвращаемое значение:
// - error - ошибка.
//
func (s *MetricService) Push(ctx context.Context, name, kind, value string) error {
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
		if v, err := s.storage.Get(ctx, name); err == nil {
			if oldVal, ok := v.GetValue().(metric.Counter); ok {
				val += oldVal
			}
		}
		record.SetValue(val)
	}

	return s.storage.Push(ctx, name, record)
}

// PushCounter добавляет метрику типа counter с указанным именем и значением.
//
// Параметры:
// - ctx - контекст;
// - name - имя метрики;
// - value - значение метрики.
//
// Возвращаемое значение:
// - error - ошибка.
//
func (s *MetricService) PushCounter(ctx context.Context, name string, value metric.Counter) (metric.Counter, error) {
	record, err := storage.NewRecord(name)
	if err != nil {
		return 0, err
	}

	if v, err := s.storage.Get(ctx, name); err == nil {
		if oldVal, ok := v.GetValue().(metric.Counter); ok {
			value += oldVal
		}
	}
	record.SetValue(value)

	return value, s.storage.Push(ctx, name, record)
}

// PushGauge добавляет метрику типа gauge с указанным именем и значением.
//
// Параметры:
// - ctx - контекст;
// - name - имя метрики;
// - value - значение метрики.
//
// Возвращаемое значение:
// - error - ошибка.
//
func (s *MetricService) PushGauge(ctx context.Context, name string, value metric.Gauge) (metric.Gauge, error) {
	record, err := storage.NewRecord(name)
	if err != nil {
		return 0, err
	}

	record.SetValue(value)

	err = s.storage.Push(ctx, name, record)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// PushBatch добавляет набор метрик.
//
// Параметры:
// - ctx - контекст;
// - records - набор метрик.
//
// Возвращаемое значение:
// - error - ошибка.
//
func (s MetricService) PushBatch(ctx context.Context, records []storage.Record) error {
	data := make([]storage.Record, 0, len(records))
	cache := make(map[string]int)
	counters := make(map[string]metric.Counter)

	for _, v := range records {
		id := v.GetName()
		value := v.GetValue()

		if i, ok := cache[id]; ok {
			if value.IsCounter() {
				if oldValue, ok := counters[id]; ok {
					value = oldValue + value.(metric.Counter)
				}
				counters[id] = value.(metric.Counter)
			}

			data[i].SetValue(value)

			continue
		}

		record, err := storage.NewRecord(id)
		if err != nil {
			return err
		}

		if value.IsCounter() {
			storRecord, err := s.Get(ctx, id, value.Kind())
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

// Get получает текущее значение метрики с указанным именем и типом.
//
// Параметры:
// - ctx - контекст;
// - name - имя метрики;
// - kind - тип метрики.
//
// Возвращаемое значение:
// - *storage.Record - текущее значение метрики;
// - error - ошибка.
//
func (s MetricService) Get(ctx context.Context, name, kind string) (*storage.Record, error) {
	if _, err := metric.GetKind(kind); err != nil {
		return nil, err
	}

	record, err := s.storage.Get(ctx, name)
	if err != nil {
		return nil, metric.ErrorMetricNotFound
	}

	return record, nil
}

// GetAll получает текущее значение всех метрик.
//
// Параметры:
// - ctx - контекст.
//
// Возвращаемое значение:
// - []storage.Record - текущее значение всех метрик.
//
func (s MetricService) GetAll(ctx context.Context) []storage.Record {
	records, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil
	}

	return records
}

// Ping проверяет состояние хранилища метрик.
//
// Параметры:
// - ctx - контекст.
//
// Возвращаемое значение:
// - error - ошибка.
//
func (s MetricService) Ping(ctx context.Context) error {
	dbStorage, ok := s.storage.(StorageWithPing)
	if !ok {
		return ErrNotSupportedMethod
	}

	return dbStorage.Ping(ctx)
}
