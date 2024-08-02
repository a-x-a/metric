package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type (
	// ServerConfig - настройки сервера.
	ServerConfig struct {
		// ListenAddress - адрес сервера сбора метрик
		ListenAddress string `env:"ADDRESS" json:"address"`
		// StoreInterval - интервал времени в секундах, по истечении которого
		// текущие показания сервера сохраняются на диск
		// (по умолчанию 300 секунд, значение `0` делает запись синхронной).
		StoreInterval time.Duration `env:"STORE_INTERVAL" json:"store_interval"`
		// FileStoregePath - полное имя файла, куда сохраняются текущие значения
		// (по умолчанию `/tmp/metrics-db.json`, пустое значение отключает функцию записи на диск).
		FileStoregePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
		// RestoreOnStart - булево значение (`true/false`),
		// определяющее, загружать или нет ранее сохранённые значения
		// из указанного файла при старте сервера (по умолчанию `true`).
		Restore bool `env:"RESTORE" json:"restore"`
		// DatabaseDSN - строка с адресом подключения к БД.
		DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`
		// Key - ключ подписи
		Key string `env:"KEY" json:"key"`
		// CryptoKey - путь до файла с приватным ключом
		CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
		// TrustedSubnet - доверенная подсеть, строковое представление бесклассовой адресации (CIDR)
		TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	}
)

// NewServerConfig - создаёт экземпляр настроек сервера.
func NewServerConfig() ServerConfig {
	storeInterval := 300
	cfg := ServerConfig{
		ListenAddress:   "localhost:8080",
		FileStoregePath: "/tmp/metrics-db.json",
		Restore:         true,
		DatabaseDSN:     "",
		Key:             "",
		CryptoKey:       "",
		TrustedSubnet:   "",
		StoreInterval:   time.Duration(storeInterval) * time.Second,
	}

	return cfg
}

func (cfg *ServerConfig) UnmarshalJSON(b []byte) error {
	type Alias ServerConfig

	var tmp struct {
		StoreInterval string `json:"store_interval"`
		*Alias
	}

	tmp.Alias = (*Alias)(cfg)

	var err error
	if err = json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	if len(tmp.StoreInterval) != 0 {
		cfg.StoreInterval, err = time.ParseDuration(tmp.StoreInterval)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg *ServerConfig) Parse() error {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	tmp := NewServerConfig()

	declarateServerFlags(&tmp)

	storeInterval := 300
	if flag.Lookup("i") == nil {
		flag.IntVar(&storeInterval, "i", storeInterval, "интервал сохранения текущих показаний сервера на диск")
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
			cfg.ListenAddress = tmp.ListenAddress
		case "i":
			cfg.StoreInterval = time.Duration(storeInterval) * time.Second
		case "f":
			cfg.FileStoregePath = tmp.FileStoregePath
		case "r":
			cfg.Restore = tmp.Restore
		case "d":
			cfg.DatabaseDSN = tmp.DatabaseDSN
		case "k":
			cfg.Key = tmp.Key
		case "crypto-key":
			cfg.CryptoKey = tmp.CryptoKey
		case "t":
			cfg.TrustedSubnet = tmp.TrustedSubnet
		}
	})

	if err := env.Parse(cfg); err != nil {
		return err
	}

	return nil
}

func declarateServerFlags(cfg *ServerConfig) {
	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ListenAddress, "a", cfg.ListenAddress, "адрес и порт сервера сбора метрик")
	}

	if flag.Lookup("f") == nil {
		flag.StringVar(&cfg.FileStoregePath, "f", cfg.FileStoregePath, "полное имя файла, куда сохраняются текущие значения")
	}

	if flag.Lookup("r") == nil {
		flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "загружать или нет ранее сохранённые значения из файла при старте")
	}

	if flag.Lookup("d") == nil {
		flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "строка с адресом подключения к БД")
	}

	if flag.Lookup("k") == nil {
		flag.StringVar(&cfg.Key, "k", cfg.Key, "ключ подписи")
	}

	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "путь до файла с приватным ключом")
	}

	if flag.Lookup("t") == nil {
		flag.StringVar(&cfg.TrustedSubnet, "t", cfg.TrustedSubnet, "доверенная подсеть в нотации CIDR")
	}
}
