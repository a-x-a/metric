package storage

import "github.com/a-x-a/go-metric/internal/models/metric"

type (
	Record struct {
		name  string
		value metric.Metric
	}

	Storage interface {
		Save(name string, rec Record) error
	}
)

func NewRecord(name string) Record {
	return Record{name: name}
}

func (rec *Record) SetValue(value metric.Metric) {
	rec.value = value
}
