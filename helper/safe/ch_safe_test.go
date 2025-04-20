package safe

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Mock struct {
}

func TestLock_Performance(t *testing.T) {
	fmt.Printf("cost1: %d ns\n", cost())
	fmt.Printf("cost2: %d ns\n", cost())
}

func cost() int64 {
	st := time.Now()
	ch := NewSafeChannel[Mock](3)
	wg1 := sync.WaitGroup{}
	wg1.Add(2)
	go ch.Range(func(mock Mock) bool {
		return true
	})
	go func() {
		defer wg1.Done()
		for i := 0; i < 10000; i++ {
			ch.Send(Mock{})
		}
	}()
	go func() {
		defer wg1.Done()
		for i := 0; i < 10000; i++ {
			ch.Send(Mock{})
		}
	}()
	wg1.Wait()
	dur := time.Since(st).Nanoseconds()
	return dur
}
