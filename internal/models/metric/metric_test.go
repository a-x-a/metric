package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "new metrics create",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, NewMetrics())
		})
	}
}

func TestMetrics_Poll(t *testing.T) {
	metric := NewMetrics()
	tests := []struct {
		name  string
		m     *Metrics
		count Counter
	}{
		{
			name:  "poll 1",
			m:     metric,
			count: Counter(1),
		},
		{
			name:  "poll 2",
			m:     metric,
			count: Counter(2),
		},
		{
			name:  "poll 3",
			m:     metric,
			count: Counter(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Poll()
			assert.Equal(t, tt.count, tt.m.PollCount)
		})
	}
}
