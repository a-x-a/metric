package main

import "os"

func main() {
	go os.Exit(1) // want "запрещено использовать прямой вызов os.Exit в функции main пакета main"
}
