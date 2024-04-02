package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounter_Kind(t *testing.T) {
	tests := []struct {
		name string
		c    Counter
		want string
	}{
		{
			name: "counter kind",
			c:    Counter(123),
			want: "counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.Kind())
		})
	}
}

func TestCounter_String(t *testing.T) {
	tests := []struct {
		name string
		c    Counter
		want string
	}{
		{
			name: "counter to string",
			c:    Counter(123),
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.String())
		})
	}
}

func TestCounter_IsCounter(t *testing.T) {
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

func TestCounter_IsGauge(t *testing.T) {
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
func TestToCounter(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		valid    bool
		expected Counter
	}{
		{
			name:     "positive integer",
			value:    "13",
			valid:    true,
			expected: 13,
		},
		{
			name:     "zero integer",
			value:    "0",
			valid:    true,
			expected: 0,
		},
		{
			name:     "negative integer",
			value:    "-13",
			valid:    true,
			expected: -13,
		},
		{
			name:  "positive float",
			value: "2345678.000000",
			valid: false,
		},
		{
			name:  "zero float",
			value: "0.000000",
			valid: false,
		},
		{
			name:  "negative float",
			value: "-2345678.000000",
			valid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			metric, err := ToCounter(tc.value)
			if tc.valid {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, metric)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
