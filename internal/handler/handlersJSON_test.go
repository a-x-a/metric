package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/adapter"
)

func sendTestRequest(t *testing.T, method, path string, data []byte) *http.Response {
	rt := NewRouter(mockService{}, zap.NewNop())
	srv := httptest.NewServer(rt)
	defer srv.Close()

	body := bytes.NewReader(data)

	req, err := http.NewRequest(method, srv.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func TestUpdateJSONMetric(t *testing.T) {
	type result struct {
		code int
	}

	tt := []struct {
		name     string
		req      adapter.RequestMetric
		expected result
	}{
		{
			name: "push counter",
			req:  adapter.NewUpdateRequestMetricCounter("PollCount", 10),
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name: "push gauge",
			req:  adapter.NewUpdateRequestMetricGauge("Alloc", 13.123),
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name: "push unknown metric kind",
			req: adapter.RequestMetric{
				ID:    "X",
				MType: "unknown",
			},
			expected: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "push blank name counter",
			req:  adapter.NewUpdateRequestMetricCounter("", 10),
			expected: result{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "push blank name gauge",
			req:  adapter.NewUpdateRequestMetricGauge("", 13.123),
			expected: result{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			data, err := json.Marshal(tc.req)
			require.NoError(err)

			resp := sendTestRequest(t, http.MethodPost, "/update/", data)

			assert.Equal(tc.expected.code, resp.StatusCode)

			if tc.expected.code == http.StatusOK {
				assert.Equal("application/json", resp.Header.Get("Content-Type"))

				respBody, err := io.ReadAll(resp.Body)
				require.NoError(err)
				defer resp.Body.Close()

				var resp adapter.RequestMetric
				err = json.Unmarshal(respBody, &resp)
				require.NoError(err)

				assert.Equal(tc.req, resp)
			}
		})
	}

	t.Run("error decode data", func(t *testing.T) {
		require := require.New(t)

		data := []byte(`invalid`)
		resp := sendTestRequest(t, http.MethodPost, "/update/", data)

		_, err := io.ReadAll(resp.Body)
		require.NoError(err)
		defer resp.Body.Close()

		require.Equal(http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetJSONMetric(t *testing.T) {
	type result struct {
		code int
		body adapter.RequestMetric
	}

	tt := []struct {
		name     string
		req      adapter.RequestMetric
		expected result
	}{
		{
			name: "get counter",
			req:  adapter.NewGetRequestMetricCounter("PollCount"),
			expected: result{
				code: http.StatusOK,
				body: adapter.NewUpdateRequestMetricCounter("PollCount", 123),
			},
		},
		{
			name: "get gauge",
			req:  adapter.NewGetRequestMetricGauge("Alloc"),
			expected: result{
				code: http.StatusOK,
				body: adapter.NewUpdateRequestMetricGauge("Alloc", 12.345),
			},
		},
		{
			name: "get unknown metric kind",
			req:  adapter.RequestMetric{ID: "Alloc", MType: "unknown"},
			expected: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "get unknown counter",
			req:  adapter.NewGetRequestMetricCounter("unknown"),
			expected: result{
				code: http.StatusNotFound,
			},
		},
		{
			name: "get unknown gauge",
			req:  adapter.NewGetRequestMetricGauge("unknown"),
			expected: result{
				code: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			data, err := json.Marshal(tc.req)
			require.NoError(err)

			resp := sendTestRequest(t, http.MethodPost, "/value/", data)

			assert.Equal(tc.expected.code, resp.StatusCode)

			if tc.expected.code == http.StatusOK {
				assert.Equal("application/json", resp.Header.Get("Content-Type"))

				respBody, err := io.ReadAll(resp.Body)
				require.NoError(err)
				defer resp.Body.Close()

				var resp adapter.RequestMetric
				err = json.Unmarshal(respBody, &resp)
				require.NoError(err)

				assert.Equal(tc.expected.body, resp)
			}
		})
	}

	t.Run("error decode data", func(t *testing.T) {
		require := require.New(t)

		data := []byte(`invalid`)
		resp := sendTestRequest(t, http.MethodPost, "/value/", data)

		_, err := io.ReadAll(resp.Body)
		require.NoError(err)
		defer resp.Body.Close()

		require.Equal(http.StatusBadRequest, resp.StatusCode)
	})
}
