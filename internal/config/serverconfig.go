package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type (
	ServerConfig struct {
		// ListenAddress - адрес сервера сбора метрик
		ListenAddress string `env:"ADDRESS"`
		// StoreInterval - интервал времени в секундах, по истечении которого
		// текущие показания сервера сохраняются на диск
		// (по умолчанию 300 секунд, значение `0` делает запись синхронной).
		StoreInterval time.Duration `env:"STORE_INTERVAL"`
		// FileStoregePath - полное имя файла, куда сохраняются текущие значения
		// (по умолчанию `/tmp/metrics-db.json`, пустое значение отключает функцию записи на диск).
		FileStoregePath string `env:"FILE_STORAGE_PATH"`
		// RestoreOnStart - булево значение (`true/false`),
		// определяющее, загружать или нет ранее сохранённые значения
		// из указанного файла при старте сервера (по умолчанию `true`).
		Restore bool `env:"RESTORE"`
		// DatabaseDSN - строка с адресом подключения к БД.
		DatabaseDSN string `env:"DATABASE_DSN"`
	}
)

func NewServerConfig() ServerConfig {
	storeInterval := 300
	cfg := ServerConfig{
		ListenAddress:   "localhost:8080",
		FileStoregePath: "/tmp/metrics-db.json",
		Restore:         true,
		DatabaseDSN:     "",
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ListenAddress, "a", cfg.ListenAddress, "адрес и порт сервера сбора метрик")
	}

	if flag.Lookup("i") == nil {
		flag.IntVar(&storeInterval, "i", storeInterval, "интервал сохранения текущих показаний сервера на диск")
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

	flag.Parse()

	cfg.StoreInterval = time.Duration(storeInterval) * time.Second

	_ = env.Parse(&cfg)

	return cfg
}
