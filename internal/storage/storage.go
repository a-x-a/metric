package storage

import "github.com/a-x-a/go-metric/internal/models/metric"

type (
	Record struct {
		name  string
		value metric.Metric
	}

	Storage interface {
		Push(name string, record Record) error
		Get(name string) (Record, bool)
		GetAll() []Record
	}
)

func NewRecord(name string) Record {
	return Record{name: name}
}

func (r *Record) SetValue(value metric.Metric) {
	r.value = value
}

func (r *Record) GetValue() metric.Metric {
	return r.value
}

func (r *Record) GetName() string {
	return r.name
}
