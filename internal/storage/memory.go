package storage

import (
	"context"
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

func (m *memStorage) Push(ctx context.Context, name string, record Record) error {
	m.Lock()
	defer m.Unlock()
	m.data[name] = record

	return nil
}

func (m *memStorage) Get(ctx context.Context, name string) (*Record, error) {
	m.Lock()
	defer m.Unlock()

	record, ok := m.data[name]
	if !ok {
		return nil, ErrNotFound
	}

	return &record, nil
}

func (m *memStorage) GetAll(ctx context.Context) ([]Record, error) {
	records := make([]Record, len(m.data))
	i := 0

	m.Lock()
	defer m.Unlock()
	for _, v := range m.data {
		records[i] = v
		i++
	}

	return records, nil
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
