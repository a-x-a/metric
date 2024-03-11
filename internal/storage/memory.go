package storage

import "sync"

type memStorage struct {
	data map[string]Record
	sync.Mutex
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
