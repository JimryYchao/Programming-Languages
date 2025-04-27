package gostd

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/*
! AddT: AddInt32, AddInt64, AddUint32, AddUint64, AddUintptr
! SwapT: SwapInt32, SwapInt64, SwapPointer, SwapUint32, SwapUint64, SwapUintptr
! CompareAndSwapT: CompareAndSwapInt32, CompareAndSwapInt64, CompareAndSwapPointer, CompareAndSwapUint32, CompareAndSwapUint64, CompareAndSwapUintptr
! LoadT: LoadInt32, LoadInt64, LoadPointer, LoadUint32, LoadUint64, LoadUintptr
! StoreT: StoreInt32, StoreInt64, StorePointer, StoreUint32, StoreUint64, StoreUintptr
! atomic.Bool, Int32, Pointer, Uint32, Uint64, Uintptr, Value
*/

func TestAtomicValue(t *testing.T) {
	var config atomic.Value // holds current server configuration
	// Create initial config value and store into config.
	config.Store(loadConfig())
	done := make(chan bool)
	go func() {
		// Reload config every 10 seconds
		// and update config value with the new version.
		for {
			select {
			case <-done:
				return
			default:
				time.Sleep(1 * time.Second)
				config.Store(loadConfig())
				log(">>>> update config <<<<")
			}
		}
	}()

	var wg sync.WaitGroup
	// Create worker goroutines that handle incoming requests
	// using the latest config value.
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-time.After(time.Second):
					r := <-requests()
					c := config.Load()
					log("load config")
					handle(r, c)
				case <-done:
					return
				}
			}
			// Handle request r using config c.
		}()
	}
	time.Sleep(3 * time.Second)
	close(done)
	wg.Wait()

}

func loadConfig() map[string]string {
	return make(map[string]string)
}

func requests() chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		c <- 1
	}()
	return c
}

func handle(_ int, _ any) {
	log("handle request")
}
