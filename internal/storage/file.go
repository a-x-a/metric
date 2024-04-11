package storage

import (
	"encoding/json"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

type withFileStorage struct {
	*memStorage
	sync.Mutex
	path     string
	syncMode bool
	logger   *zap.Logger
}

var _ Storage = &withFileStorage{}

func NewWithFileStorage(path string, syncMode bool, log *zap.Logger) *withFileStorage {
	return &withFileStorage{
		memStorage: NewMemStorage(),
		path:       path,
		syncMode:   syncMode,
		logger:     log,
	}
}

func (m *withFileStorage) Push(name string, record Record) error {
	if err := m.memStorage.Push(name, record); err != nil {
		return err
	}

	if m.syncMode {
		return m.Save()
	}

	return nil
}

func (m *withFileStorage) Save() error {
	m.Lock()
	defer m.Unlock()

	m.logger.Info("start save storage to file", zap.String("file", m.path))

	f, err := os.OpenFile(m.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	snapshot := m.memStorage.GetSnapShot()

	if err := encoder.Encode(snapshot.data); err != nil {
		return err
	}

	m.logger.Info("saved storage to file", zap.String("file", m.path))

	return nil
}

func (m *withFileStorage) Load() error {
	m.Lock()
	defer m.Unlock()

	m.logger.Info("loading storage from file", zap.String("file", m.path))

	file, err := os.Open(m.path)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&m.data); err != nil {
		return err
	}

	m.logger.Info("storage loded from file", zap.String("file", m.path))

	return nil
}

func (m *withFileStorage) Close() error {
	return m.Save()
}

type JSONMetric struct {
	Name  string  `json:"name"`            // имя метрики
	Kind  string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func recordToJSONMetric(r Record) JSONMetric {
	j := JSONMetric{
		Name: r.name,
		Kind: r.value.Kind(),
	}

	switch {
	case r.value.IsCounter():
		if v, ok := r.GetValue().(metric.Counter); ok {
			j.Delta = int64(v)
		}
	case r.value.IsGauge():
		if v, ok := r.GetValue().(metric.Gauge); ok {
			j.Value = float64(v)
		}
	}

	return j
}

func jsonMetricToRecord(j JSONMetric, r *Record) {
	r.name = j.Name
	kind, _ := metric.GetKind(j.Kind)

	switch kind {
	case metric.KindGauge:
		val := metric.Gauge(j.Value)
		r.SetValue(val)
	case metric.KindCounter:
		val := metric.Counter(j.Delta)
		r.SetValue(val)
	}
}
