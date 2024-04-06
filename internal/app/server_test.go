package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/storage"
)

func TestNewServer(t *testing.T) {
	require := require.New(t)

	t.Run("create new server", func(t *testing.T) {
		got := NewServer()
		require.NotNil(got)
	})
}

func Test_serverRunWithMemStorage(t *testing.T) {
	require := require.New(t)

	stor := storage.NewMemStorage()
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run(ctx)
	}()

	time.Sleep(time.Second * 1)

	conn, err := net.Dial("tcp", srv.httpServer.Addr)
	require.NoError(err)
	defer conn.Close()
	require.NotNil(conn)

	srv.Shutdown(syscall.SIGTERM)

	wg.Wait()
}

func Test_serverRunWithFileStorage(t *testing.T) {
	require := require.New(t)

	fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"
	stor := storage.NewWithFileStorage(fileName, false)
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run(ctx)
	}()

	time.Sleep(time.Second * 1)

	conn, err := net.Dial("tcp", srv.httpServer.Addr)
	require.NoError(err)
	defer conn.Close()
	require.NotNil(conn)

	srv.Shutdown(syscall.SIGTERM)

	wg.Wait()
}

func Test_serverPanic(t *testing.T) {
	require := require.New(t)

	defer func() {
		r := recover()
		require.NotNil(r)
	}()

	stor := storage.NewMemStorage()
	cfg := config.NewServerConfig()
	cfg.LogLevel = "unknown"
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}
	ctx := context.Background()

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	srv.Run(ctx)
	// }()
}
