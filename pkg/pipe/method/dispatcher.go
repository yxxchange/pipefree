package method

import (
	"github.com/yxxchange/pipefree/pkg/bridge/pool"
	"github.com/yxxchange/pipefree/pkg/pipe/data"
	"github.com/yxxchange/pipefree/pkg/view"
)

type Dispatcher struct{}

func (d *Dispatcher) Dispatch(eventType view.EventType, node data.Node) {
	pool.GetPool().Consume(WrapNodeToEvent(eventType, node))
}
