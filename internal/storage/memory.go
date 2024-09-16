package storage

import (
	"sync"
)

type memStorage struct {
	sync.Mutex
	data map[string]Record
}

var _ Storage = &memStorage{}

func NewMemStorage() *memStorage {
	return &memStorage{
		data: make(map[string]Record),
	}
}

func (m *memStorage) Push(name string, record Record) error {
	m.Lock()
	defer m.Unlock()
	m.data[name] = record

	return nil
}

func (m *memStorage) Get(name string) (Record, bool) {
	m.Lock()
	defer m.Unlock()
	record, ok := m.data[name]

	return record, ok
}

func (m *memStorage) GetAll() []Record {
	records := make([]Record, len(m.data))
	i := 0

	m.Lock()
	defer m.Unlock()
	for _, v := range m.data {
		records[i] = v
		i++
	}

	return records
}

func (m *memStorage) GetSnapShot() *memStorage {
	m.Lock()
	defer m.Unlock()

	snap := make(map[string]Record, len(m.data))

	for k, v := range m.data {
		snap[k] = v
	}

	return &memStorage{
		data: snap,
	}
}

func (m *memStorage) Close() error {
	return nil
}
