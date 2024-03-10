package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/a-x-a/go-metric/internal/models/metric"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type saver struct{}

func (s saver) Save(name, kind, value string) error {
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
func TestUpdateHandler(t *testing.T) {
	updHandler := NewUpdateHandler(saver{})
	srv := httptest.NewServer(updHandler)
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
