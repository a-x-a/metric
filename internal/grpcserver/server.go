package grpcserver

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
	"github.com/a-x-a/go-metric/pkg/grpcapi"
)

type (
	MetricServer struct {
		grpcapi.UnimplementedMetricsServer
		grpcServer    *grpc.Server
		service       MetricService
		trustedSubnet *net.IPNet
		address       string
		notify        chan error
		// config     config.ServerConfig
		// storage    storage.Storage
		// grpcServer *grpc.Server
		// logger     *zap.Logger
		// key        security.PrivateKey
	}

	// MetricService содержит описание методов сервиса сбора метрик.
	MetricService interface {
		// Get получает текущее значение метрики с указанным именем и типом.
		Get(ctx context.Context, name, kind string) (*storage.Record, error)
		// Update обновляет значение метрики.
		Update(ctx context.Context, requestMetric metric.RequestMetric) (metric.RequestMetric, error)
		// PushBatch добавляет набор метрик.
		PushBatch(ctx context.Context, records []storage.Record) error
		UpdateBatch(ctx context.Context, requestMetrics []metric.RequestMetric) error
	}
)

var _ grpcapi.MetricsServer = MetricServer{}

func New(s MetricService, address string, trustedSubnet *net.IPNet) *MetricServer {
	grpcServer := grpc.NewServer()
	srvc := MetricServer{
		grpcServer:    grpcServer,
		service:       s,
		trustedSubnet: trustedSubnet,
		address:       address,
		notify:        make(chan error, 1),
	}

	grpcapi.RegisterMetricsServer(grpcServer, srvc)

	return &srvc
}

func (s MetricServer) Start() {
	go func() {
		listen, err := net.Listen("tcp", s.address)
		if err != nil {
			s.notify <- err
			return
		}

		s.notify <- s.grpcServer.Serve(listen)

		close(s.notify)
	}()
}

func (s MetricServer) Stop() {
	s.grpcServer.GracefulStop()
}

func (s MetricServer) Notify() chan error {
	return s.notify
}

func (s MetricServer) Get(ctx context.Context, value *grpcapi.GetMetricRequest) (*grpcapi.GetMetricResponse, error) {
	if len(value.Id) == 0 || len(value.Mtype) == 0 {
		return nil, status.Error(codes.InvalidArgument, "id and mtype is required")
	}

	r, err := s.service.Get(ctx, value.Id, value.Mtype)
	if err != nil {
		return nil, status.Error(codes.NotFound, "metric not found")
	}

	data, err := recordToGRPCMetric(*r)
	if err != nil {
		return nil, err
	}

	return &grpcapi.GetMetricResponse{Metric: data}, nil
}

func (s MetricServer) Update(ctx context.Context, value *grpcapi.UpdateMetricRequest) (*grpcapi.UpdateMetricResponse, error) {
	m := value.GetMetric()
	if len(m.Id) == 0 || len(m.Mtype) == 0 {
		return nil, status.Error(codes.InvalidArgument, "id and mtype is required")
	}

	r, err := grpcMetricToRequestMetric(m)
	if err != nil {
		return nil, err
	}

	r, err = s.service.Update(ctx, r)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data, err := requestMetricToGRPCMetric(r)
	if err != nil {
		return nil, err
	}

	return &grpcapi.UpdateMetricResponse{Metric: data}, nil
}

func (s MetricServer) BatchUpdate(ctx context.Context, batch *grpcapi.BatchUpdateMetricRequest) (*grpcapi.BatchUpdateMetricResponse, error) {
	data := make([]metric.RequestMetric, len(batch.Data))
	for _, v := range batch.Data {
		value, err := grpcMetricToRequestMetric(v)
		if err != nil {
			return new(grpcapi.BatchUpdateMetricResponse), err
		}
		data = append(data, value)
	}

	s.service.UpdateBatch(ctx, data)
	return new(grpcapi.BatchUpdateMetricResponse), nil
}
