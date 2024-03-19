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
	// cfg.ListenAddress = "localhost:9092"
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}
	ctx := context.Background()
	time.AfterFunc(time.Second*10, func() {
		_ = srv.httpServer.Shutdown(ctx)
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run(ctx)
	}()

	conn, err := net.Dial("tcp", srv.httpServer.Addr)
	require.NoError(t, err)
	defer conn.Close()
	require.NotNil(t, conn)

	wg.Wait()
}
