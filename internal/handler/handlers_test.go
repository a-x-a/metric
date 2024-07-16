package handler

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/adapter"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/service/metricservice"
	"github.com/a-x-a/go-metric/internal/signer"
	"github.com/a-x-a/go-metric/internal/storage"
)

func getRecords() []storage.Record {
	records := []storage.Record{}
	record, _ := storage.NewRecord("Alloc")
	record.SetValue(metric.Gauge(12.345))
	records = append(records, record)

	record, _ = storage.NewRecord("PollCount")
	record.SetValue(metric.Counter(123))
	records = append(records, record)

	record, _ = storage.NewRecord("Random")
	record.SetValue(metric.Gauge(1313.1313))
	records = append(records, record)

	return records
}

func getRequestMetrics() []adapter.RequestMetric {
	records := []adapter.RequestMetric{}

	record := adapter.NewUpdateRequestMetricGauge("Alloc", 12.345)
	records = append(records, record)

	record = adapter.NewUpdateRequestMetricCounter("PollCount", 123)
	records = append(records, record)

	record = adapter.NewUpdateRequestMetricGauge("Random", 1313.1313)
	records = append(records, record)

	return records
}

type mockService struct{}

func (s mockService) Push(ctx context.Context, name, kind, value string) error {
	metricKind, err := metric.GetKind(kind)
	if err != nil {
		return err
	}

	switch metricKind {
	case metric.KindGauge:
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
	case metric.KindCounter:
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
	default:
		return metric.ErrorInvalidMetricKind
	}

	return nil
}

func (s mockService) PushCounter(ctx context.Context, name string, value metric.Counter) (metric.Counter, error) {
	if name == "" {
		return 0, storage.ErrInvalidName
	}

	return value, nil
}

func (s mockService) PushGauge(ctx context.Context, name string, value metric.Gauge) (metric.Gauge, error) {
	if name == "" {
		return 0, storage.ErrInvalidName
	}

	return value, nil
}

func (s mockService) Get(ctx context.Context, name, kind string) (*storage.Record, error) {
	_, err := metric.GetKind(kind)
	if err != nil {
		return nil, err
	}

	records := make(map[string]storage.Record)
	r, _ := storage.NewRecord("Alloc")
	r.SetValue(metric.Gauge(12.345))
	records["Alloc"] = r
	r, _ = storage.NewRecord("PollCount")
	r.SetValue(metric.Counter(123))
	records["PollCount"] = r
	r, _ = storage.NewRecord("Random")
	r.SetValue(metric.Gauge(1313.1313))
	records["Random"] = r

	value, ok := records[name]
	if !ok {
		return nil, metric.ErrorMetricNotFound
	}

	return &value, nil
}

func (s mockService) GetAll(ctx context.Context) []storage.Record {
	records := []storage.Record{}
	record, _ := storage.NewRecord("Alloc")
	record.SetValue(metric.Gauge(12.345))
	records = append(records, record)

	record, _ = storage.NewRecord("PollCount")
	record.SetValue(metric.Counter(123))
	records = append(records, record)

	record, _ = storage.NewRecord("Random")
	record.SetValue(metric.Gauge(1313.1313))
	records = append(records, record)

	return records
}

func (s mockService) Ping(ctx context.Context) error {
	return nil
}

func (s mockService) PushBatch(ctx context.Context, records []storage.Record) error {
	return nil
}

func TestUpdateHandler(t *testing.T) {
	rt := NewRouter(mockService{}, zap.NewNop(), "")
	srv := httptest.NewServer(rt)
	defer srv.Close()

	type result struct {
		code int
	}
	tt := []struct {
		name     string
		path     string
		method   string
		expected result
	}{
		{
			name:   "update counter",
			path:   "/update/counter/PollCount/12",
			method: http.MethodPost,
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name:   "update method not alloved",
			path:   "/update/counter/PollCount/12",
			method: http.MethodGet,
			expected: result{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "update gauge",
			path:   "/update/gauge/Sys/13.345",
			method: http.MethodPost,
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name:   "update without metric name",
			path:   "/update/gauge/12.345",
			method: http.MethodPost,
			expected: result{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "update unknown metric kind",
			path:   "/update/unknown/Sys/12.345",
			method: http.MethodPost,
			expected: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "update counter with invalid value",
			path:   "/update/counter/fail/10.0",
			method: http.MethodPost,
			expected: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "update gauge with invalid value",
			path:   "/update/gauge/fail/10.234;",
			method: http.MethodPost,
			expected: result{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := http.Client{}
			path := srv.URL + tc.path
			req, err := http.NewRequest(tc.method, path, nil)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "text/plain")
			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tc.expected.code, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestGetHandler(t *testing.T) {
	assert := assert.New(t)
	rt := NewRouter(mockService{}, zap.NewNop(), "")
	srv := httptest.NewServer(rt)
	defer srv.Close()

	type result struct {
		code int
	}
	tt := []struct {
		name     string
		path     string
		method   string
		expected result
	}{
		{
			name:   "get counter",
			path:   "/value/counter/PollCount",
			method: http.MethodGet,
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name:   "get gauge",
			path:   "/value/gauge/Alloc",
			method: http.MethodGet,
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name:   "get without metric name",
			path:   "/value/gauge/",
			method: http.MethodGet,
			expected: result{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "get unknown metric",
			path:   "/value/gauge/unknown",
			method: http.MethodGet,
			expected: result{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "get unknown metric kind",
			path:   "/value/unknown/Alloc",
			method: http.MethodGet,
			expected: result{
				code: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := http.Client{}
			path := srv.URL + tc.path
			req, err := http.NewRequest(tc.method, path, nil)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "text/plain")
			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(tc.expected.code, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestListHandler(t *testing.T) {
	assert := assert.New(t)
	log := zap.NewNop()

	ctrl := gomock.NewController(t)
	srvc := metricservice.NewMockmetricService(ctrl)

	h := newMetricHandlers(srvc, log)

	type result struct {
		code int
	}
	tt := []struct {
		name     string
		path     string
		method   string
		records  []storage.Record
		expected result
	}{
		{
			name:    "get all metrics",
			path:    "/",
			method:  http.MethodGet,
			records: nil,
			expected: result{
				code: http.StatusOK,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srvc.EXPECT().GetAll(context.Background()).Return(tc.records)

			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			h.List(w, req)

			assert.Equal(tc.expected.code, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func ExampleMetricHandlers_List() {
	log := zap.NewNop()

	ctrl := gomock.NewController(nil)
	srvc := metricservice.NewMockmetricService(ctrl)
	srvc.EXPECT().GetAll(context.Background()).Return(nil)

	h := newMetricHandlers(srvc, log)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
}

func Test_Ping(t *testing.T) {
	assert := assert.New(t)
	log := zap.NewNop()

	ctrl := gomock.NewController(nil)
	srvc := metricservice.NewMockmetricService(ctrl)

	h := newMetricHandlers(srvc, log)

	type result struct {
		code int
	}
	tt := []struct {
		name     string
		path     string
		method   string
		err      error
		expected result
	}{
		{
			name:   "ping OK",
			path:   "/ping",
			method: http.MethodGet,
			err:    nil,
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name:   "ping InternalServerError",
			path:   "/ping",
			method: http.MethodGet,
			err:    metricservice.ErrNotSupportedMethod,
			expected: result{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srvc.EXPECT().Ping(context.Background()).Return(tc.err)

			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			h.Ping(w, req)

			assert.Equal(tc.expected.code, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func ExampleMetricHandlers_Ping() {
	log := zap.NewNop()

	ctrl := gomock.NewController(nil)
	srvc := metricservice.NewMockmetricService(ctrl)
	srvc.EXPECT().Ping(context.Background()).Return(nil)

	h := newMetricHandlers(srvc, log)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	h.Ping(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
}

func TestUpdateBatchHandler(t *testing.T) {
	assert := assert.New(t)
	log := zap.NewNop()

	ctrl := gomock.NewController(t)
	srvc := metricservice.NewMockmetricService(ctrl)

	h := newMetricHandlers(srvc, log)

	sgnr := signer.New("secret")

	records := getRecords()

	data, err := json.Marshal(getRequestMetrics())
	assert.NoError(err)

	type result struct {
		code int
	}
	tt := []struct {
		name     string
		path     string
		method   string
		records  []storage.Record
		body     string
		err      error
		expected result
	}{
		{
			name:    "update batch normal",
			path:    "/updates",
			method:  http.MethodPost,
			records: records,
			body:    string(data),
			err:     nil,
			expected: result{
				code: http.StatusOK,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srvc.EXPECT().PushBatch(context.Background(), tc.records).Return(tc.err)

			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))

			hash, err := sgnr.Hash([]byte(tc.body))
			assert.NoError(err)

			req.Header.Set("HashSHA256", hex.EncodeToString(hash))

			w := httptest.NewRecorder()

			h.UpdateBatch(w, req)

			assert.Equal(tc.expected.code, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func ExampleMetricHandlers_UpdateBatch() {
	ctrl := gomock.NewController(nil)
	srvc := metricservice.NewMockmetricService(ctrl)
	srvc.EXPECT().PushBatch(context.Background(), getRecords()).Return(nil)

	log := zap.NewNop()
	h := newMetricHandlers(srvc, log)

	data, err := json.Marshal(getRequestMetrics())
	if err != nil {
		return
	}

	sgnr := signer.New("secret")
	hash, err := sgnr.Hash(data)
	if err != nil {
		return
	}

	bodyReader := strings.NewReader(string(data))

	req := httptest.NewRequest(http.MethodPost, "/updates", bodyReader)
	req.Header.Set("HashSHA256", hex.EncodeToString(hash))

	w := httptest.NewRecorder()

	h.UpdateBatch(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
}
