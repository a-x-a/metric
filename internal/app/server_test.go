package app

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/storage"
)

func TestNewServer(t *testing.T) {
	t.Run("create new server", func(t *testing.T) {
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
		srv:     &http.Server{},
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
			a:       "localhost:8088",
			wantErr: false,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*10, cancel)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.srv.Addr = tt.a

			go func() {
				tt.s.Run(ctx)
			}()

			conn, err := net.Dial("tcp", tt.a)
			require.NoError(t, err)
			defer conn.Close()

			require.NotNil(t, conn)
		})
	}

	// cancel()
}
