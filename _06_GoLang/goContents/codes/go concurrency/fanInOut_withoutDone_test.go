package examples

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"testing"
)

// // 重复一个 fn 并持续发送数据到 chan stream
// func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
// 	stream := make(chan T)
// 	go func() {
// 		defer close(stream)
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case stream <- fn():
// 			}
// 		}
// 	}()
// 	return stream
// }

// // 从 chan stream 取数据
//
//	func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
//		taken := make(chan T)
//		go func() {
//			defer close(taken)
//			for i := 0; i < n; i++ {
//				select {
//				case <-done:
//					return
//				case taken <- <-stream:
//				}
//			}
//		}()
//		return taken
//	}
//
// ! 构造多个工作例程
func fanOut[T any](n int, in <-chan T, work func(in <-chan T) <-chan T) ([]<-chan T, error) {
	if n < 1 {
		return nil, fmt.Errorf("n(%d) is less than 1.", n)
	}
	if n > runtime.NumCPU() {
		n = runtime.NumCPU()
	}
	fanOutChs := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		fanOutChs[i] = work(in)
	}
	return fanOutChs, nil
}

// ! 合并多个例程的 chan 通道为一个 chan // or name func merge
func fanIn[T any](cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	output := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			out <- i
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func withoutDone_generate(start int, n int) <-chan int {
	out := make(chan int)
	go func() {
		for range n {
			out <- start
			start++
		}
		close(out)
	}()
	return out
}

// 素数查找
func primeFinder(in <-chan int) <-chan int {
	isPrime := func(num int) bool {
		for i := num - 1; i > 1; i-- {
			if num%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		var ok bool
		var randomInt int
		defer close(primes)
		for {
			if randomInt, ok = <-in; ok {
				if isPrime(randomInt) {
					primes <- randomInt
				}
			} else {
				return
			}
		}
	}()
	return primes
}

func TestFanInOut(t *testing.T) {
	intStream := withoutDone_generate(400000000, 200)
	// fan out
	fanOutChs, err := fanOut(16, intStream, primeFinder)
	if err != nil {
		log.Fatal(err)
	}
	// fan in
	outStream := fanIn(fanOutChs...)

	for r := range outStream {
		fmt.Println(r)
	}
}
