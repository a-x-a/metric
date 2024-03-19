package app

import (
	"context"
	"net"
	"net/http"
	"sync"
	"testing"

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
		srv:     &http.Server{Addr: "localhost:9090"},
	}
	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run(ctx)
	}()

	conn, err := net.Dial("tcp", srv.srv.Addr)
	require.NoError(t, err)
	defer conn.Close()
	require.NotNil(t, conn)
	_ = srv.srv.Shutdown(ctx)

	wg.Wait()
}
