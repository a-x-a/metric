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
	agent struct {
		Config config.AgentConfig
	}
)

func NewAgent() *agent {
	return &agent{Config: config.NewAgentConfig()}
}

func (app *agent) Poll(ctx context.Context, metrics *metric.Metrics) {
	ticker := time.NewTicker(app.Config.PollInterval)
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

func (app *agent) Report(ctx context.Context, metrics *metric.Metrics) {
	ticker := time.NewTicker(app.Config.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := sender.SendMetrics(app.Config.ServerAddress, app.Config.PollInterval, *metrics)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
