package app

import (
	"context"
	"net"
	"net/http"
	"sync"
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
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}
	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run(ctx)
	}()

	time.Sleep(time.Second * 2)

	conn, err := net.Dial("tcp", srv.httpServer.Addr)
	require.NoError(t, err)
	defer conn.Close()
	require.NotNil(t, conn)

	_ = srv.httpServer.Shutdown(ctx)

	wg.Wait()
}
