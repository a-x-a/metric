package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGauge_Kind(t *testing.T) {
	tests := []struct {
		name string
		g    Gauge
		want string
	}{
		{
			name: "gauge kind",
			g:    Gauge(123.456),
			want: "gauge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.g.Kind())
		})
	}
}

func TestGauge_String(t *testing.T) {
	tests := []struct {
		name string
		g    Gauge
		want string
	}{
		{
			name: "counter to string",
			g:    Gauge(123.456),
			want: "123.456",
		},
		{
			name: "zero counter to string",
			g:    Gauge(0),
			want: "0.000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.g.String())
		})
	}
}

func TestGauge_IsCounter(t *testing.T) {
	tests := []struct {
		name string
		m    Metric
		want bool
	}{
		{
			name: "counter",
			m:    Counter(12),
			want: true,
		},
		{
			name: "gauge",
			m:    Gauge(12.34),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.IsCounter()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGauge_IsGauge(t *testing.T) {
	tests := []struct {
		name string
		m    Metric
		want bool
	}{
		{
			name: "counter",
			m:    Counter(12),
			want: false,
		},
		{
			name: "gauge",
			m:    Gauge(12.34),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.IsGauge()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToGauge(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		valid    bool
		expected Gauge
	}{
		{
			name:     "positive integer",
			value:    "13",
			valid:    true,
			expected: 13.0,
		},
		{
			name:     "zero integer",
			value:    "0",
			valid:    true,
			expected: 0.0,
		},
		{
			name:     "negative integer",
			value:    "-13",
			valid:    true,
			expected: -13.0,
		},
		{
			name:     "positive float",
			value:    "2345678.1234000",
			valid:    true,
			expected: 2345678.1234,
		},
		{
			name:     "zero float",
			value:    "0.000000",
			valid:    true,
			expected: 0.0,
		},
		{
			name:     "negative float",
			value:    "-2345678.1234000",
			valid:    true,
			expected: -2345678.1234,
		},
		{
			name:     "small positive float",
			value:    "0.654321",
			valid:    true,
			expected: 0.654321,
		},
		{
			name:  "meaningless value",
			value: "...",
			valid: false,
		},
		{
			name:  "malformed value",
			value: "0.12(",
			valid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			metric, err := ToGauge(tc.value)
			if tc.valid {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, metric)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
