package examples

import (
	"fmt"
	"testing"
)

// Pipeline
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// second stage: receive & handle, and sent to downstream
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// last stage: consumer
func TestPipeline(t *testing.T) {
	upstream := generate([]int{1, 2, 3, 4, 5, 66}...)
	for n := range sq(upstream) {
		fmt.Println(n)
	}
}
