package storage

type Storage interface {
	Save(metric, metricType, value string) error
}

type storage struct{}

func New() Storage {
	return storage{}
}

func (s storage) Save(metric, metricType, value string) error {
	return nil
}
