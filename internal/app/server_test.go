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
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/storage"
)

var log *zap.Logger = zap.NewNop()

func TestNewServer(t *testing.T) {
	require := require.New(t)

	t.Run("create new server", func(t *testing.T) {
		srv := NewServer()
		require.NotNil(srv)
	})
}

func TestNewServerWithDBOk(t *testing.T) {
	require := require.New(t)

	original, present := os.LookupEnv("DATABASE_DSN")
	os.Setenv("DATABASE_DSN", "host=localhost")
	if present {
		defer os.Setenv("DATABASE_DSN", original)
	} else {
		defer os.Unsetenv("DATABASE_DSN")
	}

	// 	defer func() {
	// 		r := recover()
	// 		require.NotNil(r)
	// 	}()

	t.Run("normal create new server with db", func(t *testing.T) {
		srv := NewServer()
		require.NotNil(srv)
	})
}

func TestNewServerWithDBError(t *testing.T) {
	require := require.New(t)

	original, present := os.LookupEnv("DATABASE_DSN")
	os.Setenv("DATABASE_DSN", "host=localhost port=port")
	if present {
		defer os.Setenv("DATABASE_DSN", original)
	} else {
		defer os.Unsetenv("DATABASE_DSN")
	}

	t.Run("panic create new server with db", func(t *testing.T) {
		defer func() {
			r := recover()
			require.NotNil(r)
		}()
		srv := NewServer()
		require.NotNil(srv)
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
		logger:     zap.NewNop(),
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

	stor := storage.NewWithFileStorage(fileName, false, log)
	cfg := config.NewServerConfig()

	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
		logger:     zap.NewNop(),
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

	stor := storage.NewMemStorage()
	cfg := config.NewServerConfig()
	cfg.FileStoregePath = ""

	srv2 := &http.Server{Addr: cfg.ListenAddress}
	defer srv2.Close()

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
		logger:     zap.NewNop(),
	}

	time.Sleep(time.Second * 2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	defer func() {
		r := recover()
		require.NotNil(r)
	}()

	srv.Run(ctx)

	srv2.Shutdown(ctx)

	wg.Wait()
}

// func Test_serverPanic(t *testing.T) {
// 	require := require.New(t)

// 	defer func() {
// 		r := recover()
// 		require.NotNil(r)
// 	}()

// 	stor := storage.NewMemStorage()
// 	cfg := config.NewServerConfig()
// 	cfg.LogLevel = "unknown"
// 	srv := server{
// 		Config:     cfg,
// 		Storage:    stor,
// 		httpServer: &http.Server{Addr: cfg.ListenAddress},
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()

// 	srv.Run(ctx)
// }

func Test_server_saveWithMemStorage(t *testing.T) {
	stor := storage.NewMemStorage()
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
		logger:     zap.NewNop(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	srv.saveStorage(ctx)
}

func Test_server_saveAndLoadWithFileStorage(t *testing.T) {
	require := require.New(t)

	f, err := os.CreateTemp(os.TempDir(), "test_3*.json")
	require.NoError(err)

	fileName := f.Name()

	err = f.Close()
	require.NoError(err)

	err = os.Remove(f.Name())
	require.NoError(err)

	stor := storage.NewWithFileStorage(fileName, true, log)
	cfg := config.NewServerConfig()
	cfg.FileStoregePath = fileName
	cfg.StoreInterval = time.Second * 2
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
		logger:     zap.NewNop(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.saveStorage(ctx)

	err = srv.loadStorage()
	require.NoError(err)
}

func Test_server_loadWitMemStorage(t *testing.T) {
	require := require.New(t)

	stor := storage.NewMemStorage()
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
		logger:     zap.NewNop(),
	}

	err := srv.loadStorage()
	require.Error(err)
	require.ErrorIs(err, ErrStorageNotSupportLoadFromFile)
}

func Test_server_loadWithError(t *testing.T) {
	var err error
	require := require.New(t)

	stor := storage.NewWithFileStorage("", true, log)
	cfg := config.NewServerConfig()
	srv := server{
		Config:     cfg,
		Storage:    stor,
		httpServer: &http.Server{Addr: cfg.ListenAddress},
		logger:     zap.NewNop(),
	}

	err = srv.loadStorage()
	require.Error(err)
}
