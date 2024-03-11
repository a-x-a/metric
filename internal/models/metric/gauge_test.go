package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
