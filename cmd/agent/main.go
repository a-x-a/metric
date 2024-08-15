package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-x-a/go-metric/internal/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	agent := app.NewAgent("info")

	ctx, cancel := context.WithCancel(context.Background())
	go agent.Run(ctx)
	fmt.Println("agent started")

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		signal := <-sigint
		cancel()
		fmt.Println("received signal: " + signal.String())
		fmt.Println("start to shutdown...")
		close(idleConnsClosed)
	}()

	<-idleConnsClosed
}
