package storage

import "time"

type (
	Storage interface {
		Push(name string, record Record) error
		Get(name string) (Record, bool)
		GetAll() []Record
	}
)

func NewDataStorage(path string, storeInterval time.Duration) Storage {
	if len(path) == 0 {
		return NewMemStorage()
	}

	return NewWithFileStorage(path, storeInterval == 0)
}
