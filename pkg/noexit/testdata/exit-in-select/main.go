package main

import (
	"os"
	"time"
)

func main() {
	signal := make(chan os.Signal)
	ticker := time.NewTicker(100)

	defer ticker.Stop()

	select {
	case <-signal:
		os.Exit(0) // want "запрещено использовать прямой вызов os.Exit в функции main пакета main"

	case <-ticker.C:
		os.Exit(1) // want "запрещено использовать прямой вызов os.Exit в функции main пакета main"
	}
}
