package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerConfig_UnmarshalJSON(t *testing.T) {
	require := require.New(t)
	tt := []struct {
		name     string
		src      string
		expected ServerConfig
	}{
		{
			name: "Load full configuration",
			src: `{
"address": "0.0.0.0:1234",
"store_interval": "5s",
"store_file": "/path/to/file.db",
"database_dsn": "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable",
"restore": true,
"key": "secret",
"crypto_key": "./path/to/private.pem"
}`,
			expected: ServerConfig{
				ListenAddress:   "0.0.0.0:1234",
				FileStoregePath: "/path/to/file.db",
				DatabaseDSN:     "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable",
				StoreInterval:   time.Duration(5) * time.Second,
				Restore:         true,
				Key:             "secret",
				CryptoKey:       "./path/to/private.pem",
			},
		},
		{
			name: "Load partial configuration",
			src:  `{"crypto_key": "./path/to/private.pem"}`,
			expected: ServerConfig{
				ListenAddress:   "localhost:8080",
				FileStoregePath: "/tmp/metrics-db.json",
				DatabaseDSN:     "",
				StoreInterval:   time.Duration(300) * time.Second,
				Restore:         true,
				Key:             "",
				CryptoKey:       "./path/to/private.pem",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewServerConfig()

			err := json.Unmarshal([]byte(tc.src), &cfg)
			require.NoError(err)

			require.Equal(tc.expected, cfg)
		})
	}
}

func TestServerConfig_UnmarshallInvalidJSON(t *testing.T) {
	require := require.New(t)
	tt := []struct {
		name string
		src  string
	}{
		{
			name: "Parse config with invalid data",
			src:  `{"address": 2}`,
		},
		{
			name: "Parse config with invalid store interval",
			src:  `{"store_interval": "_"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewServerConfig()

			err := json.Unmarshal([]byte(tc.src), &cfg)
			require.Error(err)
		})
	}
}
