package main

import "os"

func main() {
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			os.Exit(i) // want "запрещено использовать прямой вызов os.Exit в функции main пакета main"
		}
	}
}
