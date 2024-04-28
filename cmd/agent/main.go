package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-x-a/go-metric/internal/app"
)

func main() {
	agent := app.NewAgent()

	ctx, cancel := context.WithCancel(context.Background())
	go agent.Run(ctx)
	fmt.Println("agent started")
	// select {}
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		data := <-sigint
		cancel()
		fmt.Println("received signal: " + data.String())
		fmt.Println("start to shutdown...")
		close(idleConnsClosed)
	}()

	<-idleConnsClosed
}
