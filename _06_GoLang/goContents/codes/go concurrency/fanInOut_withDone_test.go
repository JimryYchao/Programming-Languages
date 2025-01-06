package examples

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// 重复一个 fn 并持续发送数据到 chan stream
func repeatFuncWithDone[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

// 从 chan stream 取数据
func takeWithDone[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()
	return taken
}

// ! 合并多个例程的 chan 通道为一个 chan // or name func merge
func fanInWithDone[T any](done <-chan int, chans ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	fannedInStream := make(chan T)
	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}
	wg.Add(len(chans))
	for _, c := range chans {
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()
	return fannedInStream
}

// fn 测试，取素数
func primeFinderWithDone(done <-chan int, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()
	return primes
}

func TestFanInOutWithDone(t *testing.T) {
	done := make(chan int)
	defer close(done)
	numFetcher := func(start int) func() int {
		return func() int {
			rt := start
			start++
			return rt
		}
	}
	start := time.Now()
	intStream := repeatFuncWithDone(done, numFetcher(10000000))

	// ! native；单例程
	// primeStream := primeFinder(done, intStream)
	// for rando := range take(done, primeStream, 50) {
	// 	fmt.Println(rando)
	// }

	// ! mul goroutines 多例程并行
	// fan out
	CPUCount := runtime.NumCPU()
	count := CPUCount
	primeFinderChans := make([]<-chan int, count)
	for i := 0; i < count; i++ {
		primeFinderChans[i] = primeFinderWithDone(done, intStream)
	}
	// fan in
	fannedInStream := fanInWithDone(done, primeFinderChans...)
	for rando := range takeWithDone(done, fannedInStream, 500) {
		fmt.Println(rando)
	}

	fmt.Print(time.Since(start))
}
