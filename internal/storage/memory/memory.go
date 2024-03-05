package memory

import (
	"strconv"
	"sync"
)

type memStorage struct {
	muGauge   sync.Mutex
	muCounter sync.Mutex
	gauge     map[string]float64
	counter   map[string]int64
}

func New() memStorage {
	return memStorage{}
}

func (s *memStorage) Save(metric, metricType, value string) error {
	switch metricType {
	case "gauge":
		err := s.saveGauge(metric, value)
		if err != nil {
			return err
		}
	case "counter":
		err := s.saveCounter(metric, value)
		if err != nil {
			return err
		}
	default:

	}

	return nil
}

func (s *memStorage) saveGauge(metric, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	s.muGauge.Lock()
	defer s.muGauge.Unlock()

	s.gauge[metric] = v

	return nil
}

func (s *memStorage) saveCounter(metric, value string) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	s.muCounter.Lock()
	defer s.muCounter.Unlock()

	s.counter[metric] += v

	return nil
}
