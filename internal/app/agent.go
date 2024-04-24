package app

import (
	"context"
	"fmt"
	"time"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/sender"
)

type (
	Agent struct {
		config config.AgentConfig
	}
)

func NewAgent() *Agent {
	return &Agent{config: config.NewAgentConfig()}
}

func (app *Agent) Poll(ctx context.Context, metrics *metric.Metrics) {
	ticker := time.NewTicker(app.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics.Poll()
		case <-ctx.Done():
			return
		}
	}
}

func (app *Agent) Report(ctx context.Context, metrics *metric.Metrics) {
	ticker := time.NewTicker(app.config.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := sender.SendMetrics(ctx, app.config.ServerAddress, app.config.PollInterval, *metrics)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
