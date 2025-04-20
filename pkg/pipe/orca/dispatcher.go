package orca

import (
	"context"
	"fmt"
	"github.com/yxxchange/pipefree/helper/safe"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"sync"
)

type EventFlow struct {
	ch *safe.Channel[model.Event]
}

func NewEventFlow(ch *safe.Channel[model.Event]) *EventFlow {
	return &EventFlow{ch: ch}
}

func (e *EventFlow) dispatch(event model.Event) {
	e.ch.Send(event)
}

type dispatcher struct {
	mu  sync.RWMutex
	dst map[model.EngineGroup][]*EventFlow

	qSize int
}

func newDispatcher() *dispatcher {
	return &dispatcher{
		dst:   make(map[model.EngineGroup][]*EventFlow),
		qSize: 1000,
	}
}

func (d *dispatcher) Register(eg model.EngineGroup) *EventFlow {
	d.mu.Lock()
	defer d.mu.Unlock()
	ch := safe.NewSafeChannel[model.Event](d.qSize)
	ef := NewEventFlow(ch)
	d.dst[eg] = append(d.dst[eg], ef)
	return ef
}

func (d *dispatcher) Dispatch(ctx context.Context, eventType model.EventType, node model.Node) error {
	done := make(chan struct{}, 1)
	safe.Go(func() {
		d.dispatch(eventType, node, done)
	})
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("dispatch timeout")
	}
}

func (d *dispatcher) dispatch(eventType model.EventType, node model.Node, done chan struct{}) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	event := wrapNodeToEvent(eventType, node)
	eg := convertNodeToEG(node)
	channels := d.findChannel(eg)
	for i := 0; i < len(channels); i++ {
		channels[i].dispatch(event)
	}
	done <- struct{}{}
}

func (d *dispatcher) findChannel(eg model.EngineGroup) []*EventFlow {
	return d.dst[eg]
}
