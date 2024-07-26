package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerConfigLoadFromFile(t *testing.T) {
	require := require.New(t)
	tt := []struct {
		name  string
		path  string
		check func(t *testing.T, cfg ServerConfig, err error)
	}{
		{
			name: "load full configuration",
			path: "testdata/full.server.config.json",
			check: func(t *testing.T, cfg ServerConfig, err error) {
				cfgExpected := ServerConfig{
					ListenAddress:   "0.0.0.0:1234",
					FileStoregePath: "/path/to/file.db",
					DatabaseDSN:     "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable",
					StoreInterval:   time.Duration(5) * time.Second,
					Restore:         true,
					Key:             "secret",
					CryptoKey:       "./path/to/private.pem",
				}

				require.NoError(err)
				require.Equal(cfgExpected, cfg)
			},
		},
		{
			name: "load partial configuration",
			path: "testdata/partial.server.config.json",
			check: func(t *testing.T, cfg ServerConfig, err error) {
				cfgDefault := NewServerConfig()
				cfgDefault.CryptoKey = "./path/to/private.pem"

				require.NoError(err)
				require.Equal(cfgDefault, cfg)
			},
		},
		{
			name: "load config with invalid json format",
			path: "testdata/invalid.server.config.1.json",
			check: func(t *testing.T, cfg ServerConfig, err error) {
				require.Error(err)
			},
		},
		{
			name: "load config with invalid store_interval",
			path: "testdata/invalid.server.config.2.json",
			check: func(t *testing.T, cfg ServerConfig, err error) {
				require.Error(err)
			},
		},
		{
			name: "load config file error",
			path: "testdata/no.file.json",
			check: func(t *testing.T, cfg ServerConfig, err error) {
				require.Error(err)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewServerConfig()
			err := loadConfigFromFile(tc.path, &cfg)

			tc.check(t, cfg, err)
		})
	}
}
func TestServerConfigParse(t *testing.T) {
	require := require.New(t)
	cfg := NewServerConfig()
	err := cfg.Parse()
	require.NoError(err)
}
