package config

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAgentConfig_UnmarshalJSON(t *testing.T) {
	require := require.New(t)
	tt := []struct {
		name     string
		src      string
		expected AgentConfig
	}{
		{
			name: "Load full configuration",
			src: `{
"address": "0.0.0.0:1234",
"poll_interval": "1s",
"report_interval": "5s",
"key": "secret",
"rate_limit": 5,
"crypto_key": "./path/to/publickey.pem"
}`,
			expected: AgentConfig{
				ServerAddress:  "0.0.0.0:1234",
				PollInterval:   1 * time.Second,
				ReportInterval: 5 * time.Second,
				RateLimit:      5,
				Key:            "secret",
				CryptoKey:      "./path/to/publickey.pem",
			},
		},
		{
			name: "Load partial configuration",
			src: `{
		"crypto_key": "./path/to/publickey.pem"
		}`,
			expected: AgentConfig{
				ServerAddress:  "localhost:8080",
				PollInterval:   2 * time.Second,
				ReportInterval: 10 * time.Second,
				RateLimit:      1,
				Key:            "",
				CryptoKey:      "./path/to/publickey.pem",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewAgentConfig()

			err := json.Unmarshal([]byte(tc.src), &cfg)
			fmt.Println("cfg=", cfg)
			require.NoError(err)
			fmt.Println("expected", tc.expected)
			require.Equal(tc.expected, cfg)
		})
	}
}

func TestAgentConfig_UnmarshallInvalidJSON(t *testing.T) {
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
			name: "Parse config with invalid poll interval",
			src:  `{"poll_interval": "_"}`,
		},
		{
			name: "Parse config with invalid report interval",
			src:  `{"report_interval": "x"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewAgentConfig()

			err := json.Unmarshal([]byte(tc.src), &cfg)
			require.Error(err)
		})
	}
}
