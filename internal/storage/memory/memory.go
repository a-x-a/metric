package memory

import (
	"strconv"
	"sync"
)

type memStorage struct {
	muGuage   sync.Mutex
	muCounter sync.Mutex
	guage     map[string]float64
	counter   map[string]int64
}

func New() memStorage {
	return memStorage{}
}

func (s *memStorage) Save(metric, metricType, value string) error {
	switch metricType {
	case "guage":
		err := s.saveGuage(metric, value)
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

func (s *memStorage) saveGuage(metric, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	s.muGuage.Lock()
	defer s.muGuage.Unlock()

	s.guage[metric] = v

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
