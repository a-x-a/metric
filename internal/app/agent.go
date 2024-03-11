package app

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/sender"
)

type (
	agentConfig struct {
		// PollInterval - частота обновления метрик, по умолчанию 2 сек
		PollInterval time.Duration
		// ReportInterval - частота отправки метрик на сервер, по умолчанию 10 сек
		ReportInterval time.Duration
		// ServerAddress - адрес сервера сбора метрик
		ServerAddress string
	}
	agent struct {
		Config agentConfig
	}
)

func newAgentConfig() agentConfig {
	cfg := agentConfig{
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		ServerAddress:  "localhost:8080",
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "адрес и порт сервера сбора метрик")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "частота обновления метрик")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "частота отправки метрик на сервер")
	flag.Parse()

	return cfg
}

func NewAgent() *agent {
	return &agent{Config: newAgentConfig()}
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
			err := sender.SendMetrics(app.Config.ServerAddress, *metrics)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
