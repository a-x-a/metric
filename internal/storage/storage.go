package storage

type (
	Storage interface {
		Push(name string, record Record) error
		Get(name string) (Record, bool)
		GetAll() []Record
	}
)
