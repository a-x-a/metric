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
