package orca

import (
	"github.com/yxxchange/pipefree/pkg/pipe/model"
)

type Dispatcher struct{}

func (d *Dispatcher) Dispatch(eventType model.EventType, node model.Node) error {
	return nil
}
