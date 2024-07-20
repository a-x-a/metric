// Package config инициализирует настройки агента.
package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type (
	// AgentConfig настройки агента.
	AgentConfig struct {
		// PollInterval частота обновления метрик, по умолчанию 2 сек
		PollInterval time.Duration `env:"POLL_INTERVAL"`
		// ReportInterval частота отправки метрик на сервер, по умолчанию 10 сек
		ReportInterval time.Duration `env:"REPORT_INTERVAL"`
		// ServerAddress адрес сервера сбора метрик
		ServerAddress string `env:"ADDRESS"`
		// Key ключ подписи
		Key string `env:"KEY"`
		// RateLimit количество одновременно исходящих запросов на сервер
		RateLimit int `env:"RATE_LIMIT"`
		//CryptoKey путь до файла с публичным ключом
		CryptoKey string `env:"CRYPTO_KEY"`
	}
)

// NewAgentConfig создаёт экземпляр настроек агента.
func NewAgentConfig() AgentConfig {
	pollInterval := 2
	reportInterval := 10
	serverAddress := "localhost:8080"
	key := ""
	rateLimit := 1
	cryptoKey := ""

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&serverAddress, "a", serverAddress, "адрес и порт сервера сбора метрик")
	}

	if flag.Lookup("p") == nil {
		flag.IntVar(&pollInterval, "p", pollInterval, "частота обновления метрик")
	}

	if flag.Lookup("r") == nil {
		flag.IntVar(&reportInterval, "r", reportInterval, "частота отправки метрик на сервер")
	}

	if flag.Lookup("k") == nil {
		flag.StringVar(&key, "k", key, "ключ подписи")
	}

	if flag.Lookup("l") == nil {
		flag.IntVar(&rateLimit, "l", rateLimit, "количество одновременно исходящих запросов на сервер")
	}

	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&cryptoKey, "crypto-key", cryptoKey, "путь до файла с публичным ключом")
	}

	flag.Parse()

	cfg := AgentConfig{
		PollInterval:   time.Duration(pollInterval) * time.Second,
		ReportInterval: time.Duration(reportInterval) * time.Second,
		ServerAddress:  serverAddress,
		Key:            key,
		RateLimit:      rateLimit,
		CryptoKey:      cryptoKey,
	}

	_ = env.Parse(&cfg)

	if cfg.RateLimit < 1 {
		cfg.RateLimit = 1
	}

	return cfg
}
