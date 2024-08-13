package grpcserver

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-x-a/go-metric/internal/logger"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
	"github.com/a-x-a/go-metric/pkg/grpcapi"
)

type (
	MetricServer struct {
		grpcapi.UnimplementedMetricsServiceServer
		grpcServer    *grpc.Server
		service       MetricService
		trustedSubnet *net.IPNet
		address       string
		notify        chan error
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

var _ grpcapi.MetricsServiceServer = MetricServer{}

func New(s MetricService, address string, trustedSubnet *net.IPNet, log *zap.Logger) *MetricServer {
	opts := make([]grpc.UnaryServerInterceptor, 0, 1)
	opts = append(opts, logger.UnaryRequestsInterceptor(log))

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))

	srvc := MetricServer{
		grpcServer:    grpcServer,
		service:       s,
		trustedSubnet: trustedSubnet,
		address:       address,
		notify:        make(chan error, 1),
	}

	grpcapi.RegisterMetricsServiceServer(grpcServer, srvc)

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

func (s MetricServer) Get(ctx context.Context, value *grpcapi.MetricsGetRequest) (*grpcapi.MetricsGetResponse, error) {
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

	return &grpcapi.MetricsGetResponse{Metric: data}, nil
}

func (s MetricServer) Update(ctx context.Context, value *grpcapi.MetricsUpdateRequest) (*grpcapi.MetricsUpdateResponse, error) {
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

	return &grpcapi.MetricsUpdateResponse{Metric: data}, nil
}

func (s MetricServer) UpdateBatch(ctx context.Context, batch *grpcapi.MetricsUpdateBatchRequest) (*grpcapi.MetricsUpdateBatchResponse, error) {
	response := new(grpcapi.MetricsUpdateBatchResponse)
	data := make([]metric.RequestMetric, len(batch.Data))
	for _, v := range batch.Data {
		value, err := grpcMetricToRequestMetric(v)
		if err != nil {
			return response, err
		}
		data = append(data, value)
	}

	s.service.UpdateBatch(ctx, data)
	return response, nil
}
