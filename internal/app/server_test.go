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
	cfg.FileStoregePath = ""
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

	time.Sleep(time.Second * 2)

	conn, err := net.Dial("tcp", srv.httpServer.Addr)
	require.NoError(err)
	defer conn.Close()
	require.NotNil(conn)

	srv.Shutdown(ctx, syscall.SIGTERM)

	wg.Wait()
}

func Test_serverRunWithFileStorage(t *testing.T) {
	require := require.New(t)

	f, err := os.CreateTemp(os.TempDir(), "test_1*.json")
	require.NoError(err)

	fileName := f.Name()

	err = f.Close()
	require.NoError(err)

	err = os.Remove(f.Name())
	require.NoError(err)

	// fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"
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

	time.Sleep(time.Second * 2)

	conn, err := net.Dial("tcp", srv.httpServer.Addr)
	require.NoError(err)
	defer conn.Close()
	require.NotNil(conn)

	srv.Shutdown(ctx, syscall.SIGTERM)

	wg.Wait()
}

func Test_serverErrorListenAndServe(t *testing.T) {
	require := require.New(t)

	defer func() {
		r := recover()
		require.NotNil(r)
	}()

	stor := storage.NewMemStorage()
	cfg := config.NewServerConfig()
	cfg.FileStoregePath = ""

	srv2 := &http.Server{Addr: cfg.ListenAddress}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = srv2.ListenAndServe()
	}()

	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	time.Sleep(time.Second * 2)

	srv.Run(ctx)

	srv2.Shutdown(ctx)

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.Run(ctx)
}

func Test_server_saveWithMemStorage(t *testing.T) {
	stor := storage.NewMemStorage()
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	srv.saveStorage(ctx)
}

func Test_server_saveWithFileStorage(t *testing.T) {
	require := require.New(t)

	f, err := os.CreateTemp(os.TempDir(), "test_3*.json")
	require.NoError(err)

	fileName := f.Name()

	err = f.Close()
	require.NoError(err)

	err = os.Remove(f.Name())
	require.NoError(err)

	// fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"
	stor := storage.NewWithFileStorage(fileName, true)
	cfg := config.NewServerConfig()
	cfg.FileStoregePath = fileName
	cfg.StoreInterval = time.Second * 2
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.saveStorage(ctx)
}

func Test_server_loadWitMemStorage(t *testing.T) {
	stor := storage.NewWithFileStorage("", true)
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	srv.loadStorage()
}

func Test_server_loadWitFileStorage(t *testing.T) {
	require := require.New(t)

	f, err := os.CreateTemp(os.TempDir(), "test_4*.json")
	require.NoError(err)

	fileName := f.Name()

	err = f.Close()
	require.NoError(err)

	err = os.Remove(f.Name())
	require.NoError(err)

	// fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"
	stor := storage.NewWithFileStorage(fileName, true)
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	srv.loadStorage()
}

func Test_server_loadFileError(t *testing.T) {
	var err error
	require := require.New(t)

	f, err := os.CreateTemp(os.TempDir(), "test_1*.json")
	require.NoError(err)

	fileName := f.Name()

	err = f.Close()
	require.NoError(err)

	defer func() {
		err = os.Remove(f.Name())
		require.NoError(err)
	}()

	stor := storage.NewWithFileStorage(fileName, true)
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
	}

	defer func() {
		r := recover()
		require.NotNil(r)
	}()

	srv.loadStorage()
}
