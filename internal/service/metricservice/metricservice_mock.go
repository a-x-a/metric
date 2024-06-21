package metricservice

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	metric "github.com/a-x-a/go-metric/internal/models/metric"
	storage "github.com/a-x-a/go-metric/internal/storage"
)

// MockmetricService is a mock of metricService interface.
type MockmetricService struct {
	ctrl     *gomock.Controller
	recorder *MockmetricServiceMockRecorder
}

// MockmetricServiceMockRecorder is the mock recorder for MockmetricService.
type MockmetricServiceMockRecorder struct {
	mock *MockmetricService
}

// NewMockmetricService creates a new mock instance.
func NewMockmetricService(ctrl *gomock.Controller) *MockmetricService {
	mock := &MockmetricService{ctrl: ctrl}
	mock.recorder = &MockmetricServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmetricService) EXPECT() *MockmetricServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockmetricService) Get(ctx context.Context, name, kind string) (*storage.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, name, kind)
	ret0, _ := ret[0].(*storage.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockmetricServiceMockRecorder) Get(ctx, name, kind any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockmetricService)(nil).Get), ctx, name, kind)
}

// GetAll mocks base method.
func (m *MockmetricService) GetAll(ctx context.Context) []storage.Record {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]storage.Record)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockmetricServiceMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockmetricService)(nil).GetAll), ctx)
}

// Ping mocks base method.
func (m *MockmetricService) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockmetricServiceMockRecorder) Ping(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockmetricService)(nil).Ping), ctx)
}

// Push mocks base method.
func (m *MockmetricService) Push(ctx context.Context, name, kind, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Push", ctx, name, kind, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push.
func (mr *MockmetricServiceMockRecorder) Push(ctx, name, kind, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockmetricService)(nil).Push), ctx, name, kind, value)
}

// PushBatch mocks base method.
func (m *MockmetricService) PushBatch(ctx context.Context, records []storage.Record) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushBatch", ctx, records)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushBatch indicates an expected call of PushBatch.
func (mr *MockmetricServiceMockRecorder) PushBatch(ctx, records any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushBatch", reflect.TypeOf((*MockmetricService)(nil).PushBatch), ctx, records)
}

// PushCounter mocks base method.
func (m *MockmetricService) PushCounter(ctx context.Context, name string, value metric.Counter) (metric.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushCounter", ctx, name, value)
	ret0, _ := ret[0].(metric.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushCounter indicates an expected call of PushCounter.
func (mr *MockmetricServiceMockRecorder) PushCounter(ctx, name, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushCounter", reflect.TypeOf((*MockmetricService)(nil).PushCounter), ctx, name, value)
}

// PushGauge mocks base method.
func (m *MockmetricService) PushGauge(ctx context.Context, name string, value metric.Gauge) (metric.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushGauge", ctx, name, value)
	ret0, _ := ret[0].(metric.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushGauge indicates an expected call of PushGauge.
func (mr *MockmetricServiceMockRecorder) PushGauge(ctx, name, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushGauge", reflect.TypeOf((*MockmetricService)(nil).PushGauge), ctx, name, value)
}
