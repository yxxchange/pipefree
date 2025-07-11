package safe

import (
	"github.com/yxxchange/pipefree/helper/log"
)

func Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("recover from panic: %v", r)
			}
		}()
		fn()
	}()
}
