package main

import "os"

func main() {
	code := 0

	if code == 0 {
		os.Exit(0) // want "запрещено использовать прямой вызов os.Exit в функции main пакета main"
	}

	os.Exit(code) // want "запрещено использовать прямой вызов os.Exit в функции main пакета main"
}
