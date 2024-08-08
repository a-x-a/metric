package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type (
	// AgentConfig настройки агента.
	AgentConfig struct {
		// PollInterval частота обновления метрик, по умолчанию 2 сек
		PollInterval time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
		// ReportInterval частота отправки метрик на сервер, по умолчанию 10 сек
		ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
		// RateLimit количество одновременно исходящих запросов на сервер
		RateLimit int `env:"RATE_LIMIT" json:"rate_limit"`
		// ServerAddress адрес сервера сбора метрик
		ServerAddress string `env:"ADDRESS" json:"address"`
		// Key ключ подписи
		Key string `env:"KEY" json:"key"`
		// CryptoKey путь до файла с публичным ключом в фомате PEM
		CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
		// Transport тип протокола используемого для передачи метрик (http или grpc)
		Transport string `env:"TRANSPORT" json:"transport"`
	}
)

// NewAgentConfig создаёт экземпляр настроек агента.
func NewAgentConfig() AgentConfig {
	return AgentConfig{
		PollInterval:   time.Duration(2) * time.Second,
		ReportInterval: time.Duration(10) * time.Second,
		RateLimit:      1,
		ServerAddress:  "localhost:8080",
		Key:            "",
		CryptoKey:      "",
		Transport:      "",
	}
}

func (cfg *AgentConfig) UnmarshalJSON(b []byte) error {
	type Alias AgentConfig

	var tmp struct {
		PollInterval   string `json:"poll_interval"`
		ReportInterval string `json:"report_interval"`
		*Alias
	}

	tmp.Alias = (*Alias)(cfg)

	var err error
	if err = json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	if len(tmp.PollInterval) != 0 {
		cfg.PollInterval, err = time.ParseDuration(tmp.PollInterval)
		if err != nil {
			return err
		}
	}

	if len(tmp.ReportInterval) != 0 {
		cfg.ReportInterval, err = time.ParseDuration(tmp.ReportInterval)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg *AgentConfig) Parse() error {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	tmp := NewAgentConfig()

	declarateAgentFlags(&tmp)

	pollInterval := 2
	if flag.Lookup("p") == nil {
		flag.IntVar(&pollInterval, "p", pollInterval, "частота обновления метрик")
	}

	reportInterval := 10
	if flag.Lookup("r") == nil {
		flag.IntVar(&reportInterval, "r", reportInterval, "частота отправки метрик на сервер")
	}

	configFile := ""
	if flag.Lookup("config") == nil && flag.Lookup("c") == nil {
		flag.StringVar(&configFile, "config", configFile, "путь до конфигурационного файла в формате JSON")
		flag.StringVar(&configFile, "c", configFile, "путь до конфигурационного файла в формате JSON (короткиф формат)")
	}

	flag.Parse()

	if len(configFile) != 0 {
		if err := loadConfigFromFile(configFile, cfg); err != nil {
			return err
		}
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			cfg.ServerAddress = tmp.ServerAddress
		case "p":
			cfg.PollInterval = time.Duration(pollInterval) * time.Second
		case "r":
			cfg.ReportInterval = time.Duration(reportInterval) * time.Second
		case "k":
			cfg.Key = tmp.Key
		case "l":
			cfg.RateLimit = tmp.RateLimit
		case "crypto-key":
			cfg.CryptoKey = tmp.CryptoKey
		}
	})

	if err := env.Parse(cfg); err != nil {
		return err
	}

	if cfg.RateLimit < 1 {
		cfg.RateLimit = 1
	}

	return nil
}

func declarateAgentFlags(cfg *AgentConfig) {
	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "адрес и порт сервера сбора метрик")
	}

	if flag.Lookup("k") == nil {
		flag.StringVar(&cfg.Key, "k", cfg.Key, "ключ подписи")
	}

	if flag.Lookup("l") == nil {
		flag.IntVar(&cfg.RateLimit, "l", cfg.RateLimit, "количество одновременно исходящих запросов на сервер")
	}

	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "путь до файла с публичным ключом в формате PEM")
	}
}
