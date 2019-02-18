package main

import (
	"os"
)

func main() {
	for i := 0; i < 10; i++ {
		dream, _ := os.Open("/dev/dream")
		dream.Close()
	}
}
