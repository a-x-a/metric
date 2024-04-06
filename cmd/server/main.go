package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-x-a/go-metric/internal/app"
)

func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	srv := app.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.Run(ctx)

	signal := <-sigint

	srv.Shutdown(ctx, signal)
}
