package metricservice

type metricStorage interface {
	Save(metric string, metricType string, value string) error
}

type metricService struct {
	storage metricStorage
}

func New(storage metricStorage) metricService {
	return metricService{storage}
}

func (s metricService) Save(metric string, metricType string, value string) error {
	return s.storage.Save(metric, metricType, value)
}
