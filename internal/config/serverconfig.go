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
		//CryptoKey - путь до файла с приватным ключом
		CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	}
)

func (cfg *ServerConfig) UnmarshalJSON(b []byte) error {
	var err error
	var tmp struct {
		ListenAddress   string `json:"address"`
		StoreInterval   string `json:"store_interval"`
		FileStoregePath string `json:"store_file"`
		Restore         bool   `json:"restore"`
		DatabaseDSN     string `json:"database_dsn"`
		Key             string `json:"key"`
		CryptoKey       string `json:"crypto_key"`
	}

	if err = json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	cfg.StoreInterval, err = time.ParseDuration(tmp.StoreInterval)
	if err != nil {
		return err
	}

	cfg.ListenAddress = tmp.ListenAddress
	cfg.FileStoregePath = tmp.FileStoregePath
	cfg.Restore = tmp.Restore
	cfg.DatabaseDSN = tmp.DatabaseDSN
	cfg.Key = tmp.Key
	cfg.CryptoKey = tmp.CryptoKey

	return nil
}

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
		StoreInterval:   time.Duration(storeInterval) * time.Second,
	}

	return cfg
}

func (cfg *ServerConfig) Parse() error {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	listenAddress := cfg.ListenAddress
	if flag.Lookup("a") == nil {
		flag.StringVar(&listenAddress, "a", listenAddress, "адрес и порт сервера сбора метрик")
	}

	storeInterval := 300
	if flag.Lookup("i") == nil {
		flag.IntVar(&storeInterval, "i", storeInterval, "интервал сохранения текущих показаний сервера на диск")
	}

	fileStoregePath := cfg.FileStoregePath
	if flag.Lookup("f") == nil {
		flag.StringVar(&fileStoregePath, "f", fileStoregePath, "полное имя файла, куда сохраняются текущие значения")
	}

	restore := cfg.Restore
	if flag.Lookup("r") == nil {
		flag.BoolVar(&restore, "r", restore, "загружать или нет ранее сохранённые значения из файла при старте")
	}

	databaseDSN := cfg.DatabaseDSN
	if flag.Lookup("d") == nil {
		flag.StringVar(&databaseDSN, "d", databaseDSN, "строка с адресом подключения к БД")
	}

	key := cfg.Key
	if flag.Lookup("k") == nil {
		flag.StringVar(&key, "k", key, "ключ подписи")
	}

	cryptoKey := cfg.CryptoKey
	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&cryptoKey, "crypto-key", cryptoKey, "путь до файла с приватным ключом")
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
			cfg.ListenAddress = listenAddress
		case "i":
			cfg.StoreInterval = time.Duration(storeInterval) * time.Second
		case "f":
			cfg.FileStoregePath = fileStoregePath
		case "r":
			cfg.Restore = restore
		case "d":
			cfg.DatabaseDSN = databaseDSN
		case "k":
			cfg.Key = key
		case "crypto-key":
			cfg.CryptoKey = cryptoKey
		}
	})

	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	return nil
}
