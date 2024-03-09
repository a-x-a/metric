package memory

import (
	"sync"

	"github.com/a-x-a/go-metric/internal/storage"
)

type memStorage struct {
	data map[string]storage.Record
	sync.RWMutex
}

var _ storage.Storage = &memStorage{}

func NewMemStorage() *memStorage {
	return &memStorage{
		data: make(map[string]storage.Record),
	}
}

func (ms *memStorage) Save(name string, rec storage.Record) error {
	ms.Lock()
	defer ms.Unlock()

	ms.data[name] = rec

	return nil
}
