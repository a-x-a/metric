package app

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/storage"
)

func TestNewServer(t *testing.T) {
	t.Run("create nw server", func(t *testing.T) {
		got := NewServer()
		require.NotNil(t, got)
	})
}

func Test_serverRun(t *testing.T) {
	stor := storage.NewMemStorage()
	cfg := config.ServerConfig{}
	srv := server{
		Config:  cfg,
		Storage: stor,
	}
	tests := []struct {
		name    string
		s       *server
		a       string
		wantErr bool
	}{
		{
			name:    "server run normal",
			s:       &srv,
			a:       "localhost:8081",
			wantErr: false,
		},
		{
			name:    "server run error",
			s:       &srv,
			a:       "localhost:8081",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv.Config.ListenAddress = tt.a

			go func() {
				err := srv.Run()
				if tt.wantErr {
					require.Error(t, err)
					return
				}
			}()

			conn, err := net.Dial("tcp", tt.a)
			require.NoError(t, err)
			defer conn.Close()

			require.NotNil(t, conn)
		})
	}
}
