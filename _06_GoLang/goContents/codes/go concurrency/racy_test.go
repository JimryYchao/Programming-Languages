package examples

import (
	"fmt"
	"sync"
	"testing"
)

func repeatFn(fn func(), n int) {
	for i := 0; i < n; i++ {
		fn()
	}
}

// racy goroutine
func processData(wg *sync.WaitGroup, rt *[]int, data int) {
	defer wg.Done()
	processedData := data * 2
	*rt = append(*rt, processedData)
}

func TestRacyGoroutine(t *testing.T) {
	repeatFn(func() {
		var wg sync.WaitGroup
		input := []int{1, 2, 3, 4, 5}
		result := []int{}
		for _, data := range input {
			wg.Add(1)
			go processData(&wg, &result, data)
		}
		wg.Wait()
		fmt.Println(result)
	}, 10)
}

// with lock
var lock sync.Mutex

func processDataUseLock(wg *sync.WaitGroup, rt *[]int, data int) {
	defer wg.Done()
	processedData := data * 2
	lock.Lock() // lock shares
	*rt = append(*rt, processedData)
	lock.Unlock()
}
func TestLockGoroutine(t *testing.T) {
	repeatFn(func() {
		var wg sync.WaitGroup
		input := []int{1, 2, 3, 4, 5}
		result := []int{}
		for _, data := range input {
			wg.Add(1)
			go processDataUseLock(&wg, &result, data)
		}
		wg.Wait()
		fmt.Println(result)
	}, 10)
}

// no share
func processDataNoShare(wg *sync.WaitGroup, rt *int, data int) {
	defer wg.Done()
	*rt = data * 2
}

func TestNoShareGoroutine(t *testing.T) {
	repeatFn(func() {
		var wg sync.WaitGroup
		input := []int{1, 2, 3, 4, 5}
		result := make([]int, len(input))
		for i, data := range input {
			wg.Add(1)
			go processDataNoShare(&wg, &result[i], data)
		}
		wg.Wait()
		fmt.Println(result)
	}, 10)
}
