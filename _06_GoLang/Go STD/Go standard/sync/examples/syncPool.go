package examples

import (
	. "sync"
)

type SyncPool struct {
	p *Pool
	l *Mutex
}

func NewSyncPool(new func() any) *SyncPool {
	return &SyncPool{&Pool{
		New: new,
	}, &Mutex{}}
}
func (p *SyncPool) SetNew(new func() any) <-chan bool {
	ok := make(chan bool)
	if p.l.TryLock() {
		p.p.New = new
		ok <- true
		close(ok)
	} else {
		new := new
		go func() {
			defer p.l.Unlock()
			defer close(ok)
			p.l.Lock()
			p.p.New = new
			ok <- true
		}()
	}
	return ok
}
func (p *SyncPool) Get() any {
	p.l.Lock()
	// log("get a pool")
	defer p.l.Unlock()
	return p.p.Get()
}
func (p *SyncPool) Put(x any) {
	p.p.Put(x)
}

func (p *SyncPool) Clear() {
	p.l.Lock()
	fn := p.p.New
	p.p.New = nil

	defer func() {
		p.p.New = fn
		p.l.Unlock()
		// log("clear pool")
	}()
	for {
		if p.p.Get() == nil {
			break
		}
	}
}
