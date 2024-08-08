package grpcserver

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
	"github.com/a-x-a/go-metric/pkg/grpcapi"
)

// recordToGRPCMetric преобразует storage.Record в grpcapi.Metric.
func recordToGRPCMetric(r storage.Record) (*grpcapi.Metric, error) {
	value := r.GetValue()
	result := &grpcapi.Metric{
		Id:    r.GetName(),
		Mtype: value.Kind(),
	}

	switch {
	case value.IsCounter():
		v, ok := value.(metric.Counter)
		if !ok {
			return nil, status.Error(codes.Internal, "fail to convert counter")
		}
		result.Delta = int64(v)
	case value.IsGauge():
		v, ok := value.(metric.Gauge)
		if !ok {
			return nil, status.Error(codes.Internal, "fail to convert gauge")
		}
		result.Value = float64(v)
	}

	return result, nil
}

// grpcMetricToRecord преобразует grpcapi.Metric в storage.Record.
func grpcMetricToRecord(mr *grpcapi.Metric) (storage.Record, error) {
	record, err := storage.NewRecord(mr.Id)
	if err != nil {
		return record, status.Error(codes.Internal, "fail to convert record")
	}

	kind, err := metric.GetKind(mr.Mtype)
	if err != nil {
		return record, status.Error(codes.Internal, "fail to get kind")
	}

	var value metric.Metric

	switch kind {
	case metric.KindCounter:
		value = metric.Counter(mr.Delta)

	case metric.KindGauge:
		value = metric.Gauge(mr.Value)
	}

	record.SetValue(value)

	return record, nil
}

// grpcMetricToRequestMetric преобразует grpcapi.Metric в metric.RequestMetric.
func grpcMetricToRequestMetric(mr *grpcapi.Metric) (metric.RequestMetric, error) {
	return metric.RequestMetric{
		ID:    mr.Id,
		MType: mr.Mtype,
		Delta: &mr.Delta,
		Value: &mr.Value,
	}, nil
}

// requestMetricToGRPCMetric преобразует grpcapi.Metric в metric.RequestMetric.
func requestMetricToGRPCMetric(mr metric.RequestMetric) (*grpcapi.Metric, error) {
	return &grpcapi.Metric{
		Id:    mr.ID,
		Mtype: mr.MType,
		Delta: *mr.Delta,
		Value: *mr.Value,
	}, nil
}
