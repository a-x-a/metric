// Package config инициализирует настройки приложения.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// loadConfigFromFile загружает конфигурацию из JSON файла.
func loadConfigFromFile[T any](path string, cfg *T) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config.loadConfigFromFile os.ReadFile: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("config.loadConfigFromFile json.Unmarshal: %w", err)
	}

	return nil
}
