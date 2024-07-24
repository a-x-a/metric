package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAgentConfigLoadFromFile(t *testing.T) {
	require := require.New(t)
	tt := []struct {
		name  string
		path  string
		check func(t *testing.T, cfg AgentConfig, err error)
	}{
		{
			name: "load full configuration",
			path: "testdata/full.agent.config.json",
			check: func(t *testing.T, cfg AgentConfig, err error) {
				cfgExpected := AgentConfig{
					ServerAddress:  "0.0.0.0:1234",
					PollInterval:   time.Duration(1) * time.Second,
					ReportInterval: time.Duration(5) * time.Second,
					RateLimit:      5,
					Key:            "secret",
					CryptoKey:      "./path/to/publickey.pem",
				}

				require.NoError(err)
				require.Equal(cfgExpected, cfg)
			},
		},
		{
			name: "load partial configuration",
			path: "testdata/partial.agent.config.json",
			check: func(t *testing.T, cfg AgentConfig, err error) {
				cfgDefault := NewAgentConfig()
				cfgDefault.CryptoKey = "./path/to/publickey.pem"

				require.NoError(err)
				require.Equal(cfgDefault, cfg)
			},
		},
		{
			name: "load config with invalid poll_interval",
			path: "testdata/invalid.agent.config.1.json",
			check: func(t *testing.T, cfg AgentConfig, err error) {
				require.Error(err)
			},
		},
		{
			name: "load config with invalid report_interval",
			path: "testdata/invalid.agent.config.2.json",
			check: func(t *testing.T, cfg AgentConfig, err error) {
				require.Error(err)
			},
		},
		{
			name: "load config with invalid json format",
			path: "testdata/invalid.agent.config.3.json",
			check: func(t *testing.T, cfg AgentConfig, err error) {
				require.Error(err)
			},
		},
		{
			name: "load config file error",
			path: "testdata/no.file.json",
			check: func(t *testing.T, cfg AgentConfig, err error) {
				require.Error(err)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewAgentConfig()
			err := loadConfigFromFile(tc.path, &cfg)

			tc.check(t, cfg, err)
		})
	}
}

func TestAgentConfigParse(t *testing.T) {
	require := require.New(t)
	cfg := NewAgentConfig()
	err := cfg.Parse()
	require.NoError(err)
}
