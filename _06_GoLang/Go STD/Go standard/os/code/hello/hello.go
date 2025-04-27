package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	done := time.After(10 * time.Second)

	for {
		select {
		case <-done:
			os.Exit(0)
		default:
			fmt.Println("Hello World")
			time.Sleep(250 * time.Millisecond)
		}
	}
}
