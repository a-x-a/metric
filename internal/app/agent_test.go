package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

// func TestNewAgent(t *testing.T) {
// 	t.Run("create nw agent", func(t *testing.T) {
// 		got := NewAgent()
// 		require.NotNil(t, got)
// 	})
// }

func Test_agent_Poll(t *testing.T) {
	cfg := config.AgentConfig{
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		ServerAddress:  "localhost:8080",
	}
	metrics := metric.NewMetrics()
	type args struct {
		ctx     context.Context
		metrics *metric.Metrics
	}
	tests := []struct {
		name string
		app  *agent
		args args
	}{
		{
			name: "poll",
			app:  &agent{Config: cfg},
			args: args{
				ctx:     context.Background(),
				metrics: metrics,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancellingCtx, cancel := context.WithCancel(tt.args.ctx)
			time.AfterFunc(tt.app.Config.PollInterval, cancel)
			tt.app.Poll(cancellingCtx, tt.args.metrics)

			require.NotEmpty(t, tt.args.metrics)
			require.NotEmpty(t, tt.args.metrics.PollCount)
		})
	}
}
