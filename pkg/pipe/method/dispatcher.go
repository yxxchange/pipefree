package method

import (
	"github.com/yxxchange/pipefree/pkg/pipe/data"
)

type Dispatcher struct{}

func (d *Dispatcher) Dispatch(eventType data.EventType, node data.Node) error {
	return nil
}
