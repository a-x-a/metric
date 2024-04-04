package storage

import (
	"errors"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

type (
	Record struct {
		Name  string
		Value metric.Metric
	}
)

var (
	// ErrInvalidName - не корректное имя записи.
	ErrInvalidName = errors.New("record: a record has to have a valid name")
)

func NewRecord(name string) (Record, error) {
	if name == "" {
		return Record{}, ErrInvalidName
	}
	return Record{Name: name}, nil
}

func (r *Record) SetValue(value metric.Metric) {
	r.Value = value
}

func (r *Record) GetValue() metric.Metric {
	return r.Value
}

func (r *Record) GetName() string {
	return r.Name
}
