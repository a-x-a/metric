package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-x-a/go-metric/internal/app"
)

// @title Go Metric
// @description This is a Metrics Collection Server.
// @version 0.1

// @host localhost:8080
// @BasePath /

// main is the entry point of the Go program.
//
// It sets up a signal channel to handle interrupt signals and creates a new server instance.
// It creates a context with a cancel function to gracefully shutdown the server.
// It runs the server in a separate goroutine and waits for an interrupt signal.
// It creates a new context with a timeout to gracefully shutdown the server.
// It calls the Shutdown function of the server with the shutdown context and the interrupt signal.
func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	srv := app.NewServer("info")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.Run(ctx)

	signal := <-sigint

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, time.Second*5)
	defer cancelShutdown()

	srv.Shutdown(ctxShutdown, signal)
}
