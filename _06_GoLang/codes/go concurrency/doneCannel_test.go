package examples

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestOrDoneBetweenGoroutines(t *testing.T) {
	var wg sync.WaitGroup
	done := make(chan any)

	cows := make(chan any, 100)
	pigs := make(chan any, 100)
	go func() {
		for {
			select {
			case <-done:
				return
			case cows <- "moo":
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case pigs <- "oink":
			}
		}
	}()

	wg.Add(1)
	go consumeCows(&wg, done, cows)
	wg.Add(1)
	go consumePigs(&wg, done, pigs)

	time.Sleep(1 * time.Millisecond)
	close(done)
	fmt.Println("close done")
	wg.Wait()
}

func consumeCows(wg *sync.WaitGroup, done <-chan any, cows <-chan any) {
	defer wg.Done()
	for cow := range orDone(done, cows) {
		// do something
		fmt.Println(cow)
	}
	fmt.Println("quit consumeCows")
}

func consumePigs(wg *sync.WaitGroup, done <-chan any, pigs <-chan any) {
	defer wg.Done()
	for pig := range orDone(done, pigs) {
		// do something
		fmt.Println(pig)
	}
	fmt.Println("quit consumePigs")
}

func orDone(done <-chan any, c <-chan any) <-chan any {
	relayStream := make(chan any)
	go func() {
		defer close(relayStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case relayStream <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return relayStream
}
