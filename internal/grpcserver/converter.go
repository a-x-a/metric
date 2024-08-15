package grpcserver

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcapi "github.com/a-x-a/go-metric/api/proto/v1"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
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
		result.Value = &grpcapi.Metric_Counter{
			Counter: int64(v),
		}
	case value.IsGauge():
		v, ok := value.(metric.Gauge)
		if !ok {
			return nil, status.Error(codes.Internal, "fail to convert gauge")
		}
		result.Value = &grpcapi.Metric_Gauge{
			Gauge: float64(v),
		}
	}

	return result, nil
}

// grpcMetricToRecord преобразует grpcapi.Metric в storage.Record.
func grpcMetricToRecord(mr *grpcapi.Metric) (storage.Record, error) {
	record, err := storage.NewRecord(mr.GetId())
	if err != nil {
		return record, status.Error(codes.Internal, "fail to convert record")
	}

	var value metric.Metric
	switch v := mr.Value.(type) {
	case *grpcapi.Metric_Counter:
		value = metric.Counter(v.Counter)
	case *grpcapi.Metric_Gauge:
		value = metric.Gauge(v.Gauge)
	default:
		return record, status.Error(codes.Internal, "fail to convert value")
	}

	record.SetValue(value)

	return record, nil
}

// grpcMetricToRequestMetric преобразует grpcapi.Metric в metric.RequestMetric.
func grpcMetricToRequestMetric(mr *grpcapi.Metric) (metric.RequestMetric, error) {
	result := metric.RequestMetric{
		ID: mr.GetId(),
	}

	switch v := mr.Value.(type) {
	case *grpcapi.Metric_Counter:
		result.MType = string(metric.KindCounter)
		val := float64(v.Counter)
		result.Value = &val
	case *grpcapi.Metric_Gauge:
		result.MType = string(metric.KindGauge)
		val := float64(v.Gauge)
		result.Value = &val
	default:
		return result, status.Error(codes.Internal, "fail to convert value")
	}

	return result, nil
}

// requestMetricToGRPCMetric преобразует grpcapi.Metric в metric.RequestMetric.
func requestMetricToGRPCMetric(mr metric.RequestMetric) (*grpcapi.Metric, error) {
	result := &grpcapi.Metric{
		Id:    mr.ID,
		Mtype: mr.MType,
	}

	switch mr.MType {
	case string(metric.KindCounter):
		result.Value = &grpcapi.Metric_Counter{
			Counter: int64(*mr.Delta),
		}
	case string(metric.KindGauge):
		result.Value = &grpcapi.Metric_Gauge{
			Gauge: float64(*mr.Value),
		}
	default:
		return nil, status.Error(codes.Internal, "fail to convert value")
	}

	return result, nil
}
