package main

import (
	"fmt"

	"github.com/a-x-a/go-metric/internal/app"
)

func main() {
	srv := app.NewServer()
	fmt.Println(srv.Run())
}
