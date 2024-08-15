package main

import (
	"context"
	"fmt"

	"github.com/a-x-a/go-metric/internal/app"
)

//	@title			Сервис сбора метрик и алертинга.
//	@description	Сервис для сбора рантайм-метрик.
//	@version		0.1

// @host		localhost:8080
// @BasePath	/

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	srv := app.NewServer("info")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv.Run(ctx)
}
