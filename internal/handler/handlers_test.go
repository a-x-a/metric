package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

type mockService struct{}

func (s mockService) Push(name, kind, value string) error {
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

func (s mockService) PushCounter(name string, value metric.Counter) (metric.Counter, error) {
	if name == "" {
		return 0, storage.ErrInvalidName
	}

	return value, nil
}

func (s mockService) PushGauge(name string, value metric.Gauge) (metric.Gauge, error) {
	if name == "" {
		return 0, storage.ErrInvalidName
	}

	return value, nil
}

func (s mockService) Get(name, kind string) (string, error) {
	_, err := metric.GetKind(kind)
	if err != nil {
		return "", err
	}

	records := map[string]string{
		"Alloc":     fmt.Sprintf("%.3f", 12.345),
		"PollCount": fmt.Sprintf("%d", 123),
		"Random":    fmt.Sprintf("%.3f", 1313.131),
	}
	value, ok := records[name]
	if !ok {
		return "", metric.ErrorMetricNotFound
	}

	return value, nil
}

func (s mockService) GetAll() []storage.Record {
	records := []storage.Record{}
	record, _ := storage.NewRecord("Alloc")
	record.SetValue(metric.Gauge(12.3456))
	records = append(records, record)

	record, _ = storage.NewRecord("PollCount")
	record.SetValue(metric.Counter(123))
	records = append(records, record)

	record, _ = storage.NewRecord("Random")
	record.SetValue(metric.Gauge(1313.1313))
	records = append(records, record)

	return records
}

func (s mockService) Ping() error {
	return nil
}

type mockServiceWithErrorPing struct {
	mockService
}

func (s mockServiceWithErrorPing) Ping() error {
	return fmt.Errorf("no ping")
}

func TestUpdateHandler(t *testing.T) {
	rt := NewRouter(mockService{}, zap.NewNop())
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
	rt := NewRouter(mockService{}, zap.NewNop())
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

			assert.Equal(t, tc.expected.code, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestListHandler(t *testing.T) {
	rt := NewRouter(mockService{}, zap.NewNop())
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
			name:   "get all metrics",
			path:   "/",
			method: http.MethodGet,
			expected: result{
				code: http.StatusOK,
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

func TestPingHandlerOk(t *testing.T) {
	rt := NewRouter(mockService{}, zap.NewNop())
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
			name:   "ping",
			path:   "/ping/",
			method: http.MethodGet,
			expected: result{
				code: http.StatusOK,
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
func TestPingHandlerError(t *testing.T) {
	rt := NewRouter(mockServiceWithErrorPing{}, zap.NewNop())
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
			name:   "ping",
			path:   "/ping/",
			method: http.MethodGet,
			expected: result{
				code: http.StatusInternalServerError,
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
