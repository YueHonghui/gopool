package gopool

import (
	"container/list"
	"sync"
)

type Pool struct {
	New      func() interface{}
	Destroy  func(interface{})
	lock     sync.Mutex
	maxLimit int
	idlelist list.List
}

func (p *Pool) SetMaxLimit(limit int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.maxLimit = limit
}

func (p *Pool) Get() interface{} {
	p.lock.Lock()
	if p.idlelist.Len() > 0 {
		defer p.lock.Unlock()
		i := p.idlelist.Front()
		return p.idlelist.Remove(i)
	}
	p.lock.Unlock()
	return p.New()
}

func (p *Pool) Put(value interface{}) {
	p.lock.Lock()
	if p.idlelist.Len() >= p.maxLimit {
		p.lock.Unlock()
		p.Destroy(value)
		return
	}
	p.idlelist.PushBack(value)
	p.lock.Unlock()
}

func (p *Pool) Clear() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for p.idlelist.Len() > 0 {
		p.idlelist.Remove(p.idlelist.Front())
	}
}
