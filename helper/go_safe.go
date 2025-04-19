package helper

import "fmt"

func Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// TODO: log
				fmt.Printf("panic: %v", r)
			}
		}()
		fn()
	}()
}

func GoWithCh(fn func(), done chan struct{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// TODO: log
				fmt.Printf("panic: %v", r)
			}
		}()
		fn()
		done <- struct{}{}
	}()
}
