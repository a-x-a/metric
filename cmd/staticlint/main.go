package main

import (
	"github.com/a-x-a/go-metric/pkg/multichecker"
)

func main() {
	chkr := multichecker.New()
	chkr.Run()
}
