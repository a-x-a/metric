package storage

import (
	"encoding/json"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/logger"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

type withFileStorage struct {
	*memStorage
	sync.Mutex
	path     string
	syncMode bool
}

var _ Storage = &withFileStorage{}

func NewWithFileStorage(path string, syncMode bool) *withFileStorage {
	return &withFileStorage{
		memStorage: NewMemStorage(),
		path:       path,
		syncMode:   syncMode,
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

	logger.Log.Info("start save storage to file", zap.String("file", m.path))

	f, err := os.OpenFile(m.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	snapshot := m.memStorage.GetSnapShot()

	if err := encoder.Encode(snapshot.data); err != nil {
		logger.Log.Info("error of save storage to file", zap.Error(err))
		return err
	}

	logger.Log.Info("saved storage to file", zap.String("file", m.path))

	return nil
}

func (m *withFileStorage) Load() error {
	m.Lock()
	defer m.Unlock()

	logger.Log.Info("loading storage from file", zap.String("file", m.path))

	file, err := os.Open(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Log.Error("storage file not found", zap.String("file", m.path))
			return nil
		}

		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&m.data); err != nil {
		return err
	}

	logger.Log.Info("storage loded from file", zap.String("file", m.path))

	return nil
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
	switch j.Kind {
	case "gauge":
		val := metric.Gauge(j.Value)
		r.SetValue(val)
	case "counter":
		val := metric.Counter(j.Delta)
		r.SetValue(val)
	}
}

func (r Record) MarshalJSON() ([]byte, error) {
	j := recordToJSONMetric(r)
	return json.Marshal(j)
}

func (r *Record) UnmarshalJSON(data []byte) error {
	j := JSONMetric{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	jsonMetricToRecord(j, r)

	return nil
}
