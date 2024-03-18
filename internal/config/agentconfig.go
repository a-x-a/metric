package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type (
	AgentConfig struct {
		// PollInterval - частота обновления метрик, по умолчанию 2 сек
		PollInterval time.Duration `env:"POLL_INTERVAL"`
		// ReportInterval - частота отправки метрик на сервер, по умолчанию 10 сек
		ReportInterval time.Duration `env:"REPORT_INTERVAL"`
		// ServerAddress - адрес сервера сбора метрик
		ServerAddress string `env:"ADDRESS"`
	}
)

func NewAgentConfig() AgentConfig {
	cfg := AgentConfig{
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		ServerAddress:  "localhost:8080",
	}

	pollInterval := 2
	reportInterval := 10

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "адрес и порт сервера сбора метрик")
	flag.IntVar(&pollInterval, "p", pollInterval, "частота обновления метрик")
	flag.IntVar(&reportInterval, "r", reportInterval, "частота отправки метрик на сервер")

	flag.Parse()

	cfg.PollInterval = time.Duration(pollInterval) * time.Second
	cfg.ReportInterval = time.Duration(reportInterval) * time.Second

	_ = env.Parse(&cfg)

	return cfg
}
