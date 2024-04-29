package storage

import (
	"encoding/json"
	"errors"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

type (
	Record struct {
		name  string
		value metric.Metric
	}
)

var (
	// ErrInvalidName - не корректное имя записи.
	ErrInvalidName = errors.New("record: a record has to have a valid name")
	// ErrNotFound - запись не найдена.
	ErrNotFound = errors.New("record: a record not found")
)

func NewRecord(name string) (Record, error) {
	if name == "" {
		return Record{}, ErrInvalidName
	}
	return Record{name: name}, nil
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

func (r Record) MarshalJSON() ([]byte, error) {
	j := recordToJSONMetric(r)
	return json.Marshal(j)
}

func (r *Record) UnmarshalJSON(data []byte) error {
	j := JSONMetric{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	jsonMetricToRecord(j, r)

	return nil
}
