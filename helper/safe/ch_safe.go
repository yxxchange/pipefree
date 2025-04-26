package safe

import (
	"sync"
)

type Channel[T interface{}] struct {
	ch     chan T
	rw     sync.RWMutex
	mu     sync.Mutex
	closed bool
}

func NewSafeChannel[T interface{}](size int) *Channel[T] {
	return &Channel[T]{
		ch: make(chan T, size),
	}
}

func (c *Channel[T]) Chan() <-chan T {
	return c.ch
}

func (c *Channel[T]) Range(fn func(T) (interrupted bool)) {
	for data := range c.ch {
		if fn(data) {
			return
		}
	}
}

func (c *Channel[T]) Send(data T) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	if c.closed {
		return
	}
	c.ch <- data
}

func (c *Channel[T]) Close() {
	c.rw.Lock()
	defer c.rw.Unlock()
	if c.closed {
		return
	}
	close(c.ch)
	c.closed = true
}
