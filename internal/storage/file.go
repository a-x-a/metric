package storage

import (
	"encoding/json"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/logger"
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

	logger.Log.Info("saved storage to file", zap.String("file", m.path), zap.Any("JSON", snapshot.data))

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
	if err := decoder.Decode(m.memStorage); err != nil {
		return err
	}

	logger.Log.Info("storage loded from file", zap.String("file", m.path))

	return nil
}
