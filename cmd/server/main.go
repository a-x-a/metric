package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-x-a/go-metric/internal/app"
)

//	@title			Сервис сбора метрик и алертинга.
//	@description	Сервис для сбора рантайм-метрик.
//	@version		0.1

// @host		localhost:8080
// @BasePath	/

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
