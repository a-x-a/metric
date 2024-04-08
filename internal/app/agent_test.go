package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

func TestNewAgent(t *testing.T) {
	require := require.New(t)

	t.Run("create new agent", func(t *testing.T) {
		got := NewAgent()
		require.NotNil(got)
	})
}

func Test_agent_Poll(t *testing.T) {
	require := require.New(t)

	cfg := config.AgentConfig{
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		ServerAddress:  "",
	}
	metrics := &metric.Metrics{}

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
			cancellingCtx, cancel := context.WithTimeout(tt.args.ctx, tt.app.Config.PollInterval*2)
			defer cancel()

			tt.app.Poll(cancellingCtx, tt.args.metrics)

			require.NotEmpty(tt.args.metrics)
			require.NotEmpty(tt.args.metrics.PollCount)
		})
	}
}

func Test_agent_Report(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/update/counter/PollCount/12") {
			rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			rw.WriteHeader(http.StatusOK)
			return
		}
		rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		rw.WriteHeader(http.StatusNotFound)
	}))

	defer server.Close()

	cfg := config.AgentConfig{
		PollInterval:   2 * time.Second,
		ReportInterval: 2 * time.Second,
		ServerAddress:  strings.TrimPrefix(server.URL, "http://"),
	}

	metrics := metric.Metrics{
		PollCount: metric.Counter(12),
	}

	tests := []struct {
		name    string
		app     *agent
		ctx     context.Context
		metrics *metric.Metrics
	}{
		{
			name:    "report",
			app:     &agent{Config: cfg},
			ctx:     context.Background(),
			metrics: &metrics,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancellingCtx, cancel := context.WithTimeout(tt.ctx, tt.app.Config.ReportInterval)
			defer cancel()
			tt.app.Report(cancellingCtx, tt.metrics)
		})
	}
}
