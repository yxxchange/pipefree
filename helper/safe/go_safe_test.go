package safe

import (
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	Go(func() {
		panic("test panic")
	})
	time.Sleep(1 * time.Second)
}
