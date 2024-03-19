package main

import (
	"context"

	"github.com/a-x-a/go-metric/internal/app"
)

func main() {
	srv := app.NewServer()
	ctx := context.Background()
	srv.Run(ctx)
}
