package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-x-a/go-metric/internal/app"
	"github.com/a-x-a/go-metric/internal/models/metric"
)

func main() {
	agent := app.NewAgent()
	ctx := context.Background()
	metric := &metric.Metrics{}

	go agent.Poll(ctx, metric)
	go agent.Report(ctx, metric)

	fmt.Println("agent started")

	// select {}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		data := <-sigint
		fmt.Println("received signal: " + data.String())
		fmt.Println("start to shutdown...")
		close(idleConnsClosed)
	}()
	<-idleConnsClosed
}
