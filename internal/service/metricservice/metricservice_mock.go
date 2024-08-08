package metricservice

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	metric "github.com/a-x-a/go-metric/internal/models/metric"
	storage "github.com/a-x-a/go-metric/internal/storage"
)

// MockMetricService is a mock of MetricService interface.
type MockMetricService struct {
	ctrl     *gomock.Controller
	recorder *MockMetricServiceMockRecorder
}

// MockMetricServiceMockRecorder is the mock recorder for MockMetricService.
type MockMetricServiceMockRecorder struct {
	mock *MockMetricService
}

// NewMockMetricService creates a new mock instance.
func NewMockMetricService(ctrl *gomock.Controller) *MockMetricService {
	mock := &MockMetricService{ctrl: ctrl}
	mock.recorder = &MockMetricServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricService) EXPECT() *MockMetricServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockMetricService) Get(ctx context.Context, name, kind string) (*storage.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, name, kind)
	ret0, _ := ret[0].(*storage.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMetricServiceMockRecorder) Get(ctx, name, kind any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMetricService)(nil).Get), ctx, name, kind)
}

// GetAll mocks base method.
func (m *MockMetricService) GetAll(ctx context.Context) []storage.Record {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]storage.Record)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricServiceMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricService)(nil).GetAll), ctx)
}

// Ping mocks base method.
func (m *MockMetricService) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockMetricServiceMockRecorder) Ping(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockMetricService)(nil).Ping), ctx)
}

// Push mocks base method.
func (m *MockMetricService) Push(ctx context.Context, name, kind, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Push", ctx, name, kind, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push.
func (mr *MockMetricServiceMockRecorder) Push(ctx, name, kind, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockMetricService)(nil).Push), ctx, name, kind, value)
}

// PushBatch mocks base method.
func (m *MockMetricService) PushBatch(ctx context.Context, records []storage.Record) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushBatch", ctx, records)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushBatch indicates an expected call of PushBatch.
func (mr *MockMetricServiceMockRecorder) PushBatch(ctx, records any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushBatch", reflect.TypeOf((*MockMetricService)(nil).PushBatch), ctx, records)
}

// PushCounter mocks base method.
func (m *MockMetricService) PushCounter(ctx context.Context, name string, value metric.Counter) (metric.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushCounter", ctx, name, value)
	ret0, _ := ret[0].(metric.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushCounter indicates an expected call of PushCounter.
func (mr *MockMetricServiceMockRecorder) PushCounter(ctx, name, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushCounter", reflect.TypeOf((*MockMetricService)(nil).PushCounter), ctx, name, value)
}

// PushGauge mocks base method.
func (m *MockMetricService) PushGauge(ctx context.Context, name string, value metric.Gauge) (metric.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushGauge", ctx, name, value)
	ret0, _ := ret[0].(metric.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushGauge indicates an expected call of PushGauge.
func (mr *MockMetricServiceMockRecorder) PushGauge(ctx, name, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushGauge", reflect.TypeOf((*MockMetricService)(nil).PushGauge), ctx, name, value)
}

// Update mocks base method.
func (m *MockMetricService) Update(ctx context.Context, requestMetric metric.RequestMetric) (metric.RequestMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, requestMetric)
	ret0, _ := ret[0].(metric.RequestMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockMetricServiceMockRecorder) Update(ctx, requestMetric any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMetricService)(nil).Update), ctx, requestMetric)
}

// UpdateBatch mocks base method.
func (m *MockMetricService) UpdateBatch(ctx context.Context, requestMetrics []metric.RequestMetric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBatch", ctx, requestMetrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBatch indicates an expected call of UpdateBatch.
func (mr *MockMetricServiceMockRecorder) UpdateBatch(ctx, requestMetrics any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBatch", reflect.TypeOf((*MockMetricService)(nil).UpdateBatch), ctx, requestMetrics)
}
