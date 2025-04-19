package safe

import (
	"sync"
)

type Channel[T interface{}] struct {
	ch     chan T
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

func (c *Channel[T]) Range(fn func(T) bool) {
	for data := range c.ch {
		if !fn(data) {
			return
		}
	}
}

func (c *Channel[T]) Send(data T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.ch <- data
}

func (c *Channel[T]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	close(c.ch)
	c.closed = true
}
