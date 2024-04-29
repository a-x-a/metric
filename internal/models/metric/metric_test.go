package metric

import (
	"context"
	"testing"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_Poll(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	metric := &Metrics{}
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

	CPUCount, err := cpu.Counts(true)
	require.NoError(err)
	CPUutilization1 := Gauge(CPUCount)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Poll(context.TODO())
			assert.Equal(tt.count, tt.m.PollCount)
			assert.Equal(CPUutilization1, tt.m.PS.CPUutilization1)
		})
	}
}
