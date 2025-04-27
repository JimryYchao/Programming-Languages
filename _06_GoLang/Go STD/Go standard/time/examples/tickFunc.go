package examples

import "time"

type Ticker struct {
	done   chan bool
	t      *time.Ticker
	f      func()
	isDone bool
}

// Stop the Ticker
func (t *Ticker) Stop() bool {
	if !t.isDone {
		t.done <- true
		t.isDone = true
		return true
	} else {
		return false
	}
}

// reset f and d
func (t *Ticker) Reset(d time.Duration, f func()) {
	if t.isDone {
		*t = *TickFunc(d, f)
	} else {
		t.t.Reset(d)
		t.f = f
	}
}

// call f per d
func TickFunc(d time.Duration, f func()) *Ticker {
	ticker := &Ticker{make(chan bool), time.NewTicker(d), f, false}
	go func() {
		defer ticker.t.Stop()
		for range ticker.t.C {
			select {
			case <-ticker.done:
				return
			default:
				ticker.f()
			}
		}
	}()
	return ticker
}
