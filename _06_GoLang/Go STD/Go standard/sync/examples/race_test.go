package examples

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimerDead(t *testing.T) {
	tm := time.NewTimer(1)
	tm.Reset(100 * time.Millisecond)
	<-tm.C
	if !tm.Stop() {
		// <-tm.C  // deadlock and race
	}
	// toChanTimed(time.NewTimer(1), make(chan int)) deadlock
}
func toChanTimed(t *time.Timer, ch chan int) {
	t.Reset(1 * time.Second)
	defer func() {
		if !t.Stop() {
			<-t.C
		}
	}()
	select {
	case ch <- 42:
	case <-t.C:
	}
}

func TestWrongRace(t *testing.T) {
	s := SharedInt{}
	var wg sync.WaitGroup
	n := 5000
	wg.Add(n * 2)
	go func() {
		for i := range n {
			i := i
			go func() {
				defer wg.Done()
				Increment(&s, i)
			}()
		}
	}()
	go func() {
		for range n {
			go func() {
				defer wg.Done()
				Answer(&s)
			}()
		}
	}()
	wg.Wait()
	fmt.Printf("last %d\n", s.val)
}

type SharedInt struct {
	mu  sync.Mutex
	val int
}

func Answer(si *SharedInt) {
	si.mu.Lock()
	defer si.mu.Unlock() // 🔒
	_ = si.val           // 🔒
}

func (si *SharedInt) SetVal(val int) {
	si.mu.Lock()
	defer si.mu.Unlock() // 🔒
	si.val = val         // 🔒
}

func (si *SharedInt) Val() int {
	si.mu.Lock()
	defer si.mu.Unlock() // 🔒
	return si.val        // 🔒
}

/*
The race condition manifests between the call to Val and the one to SetVal.

	func Increment(si *SharedInt) {
	  v := si.Val() // 🔒
	  v++           // 😱😱😱😱😱😱
	  si.SetVal(v)  // 🔒
	}

use

	func Increment(si *SharedInt) {
		si.mu.Lock()
		defer si.mu.Unlock() // 🔒
		si.val++             // 🔒
	}
*/
func Increment(si *SharedInt, i int) {
	si.SetVal(si.Val() + 1) // Locker race
}
