package main

import "os"

func main() {
	defer os.Exit(1) // want "запрещено использовать прямой вызов os.Exit в функции main пакета main"

	defer func() {
		os.Exit(0)
	}()
}
