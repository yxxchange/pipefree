package orca

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/helper/safe"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"sync"
	"time"
)

type EventFlow struct {
	bch chan []byte
	sch *safe.Channel[model.Event]
}

func NewEventFlow(ch *safe.Channel[model.Event]) *EventFlow {
	return &EventFlow{
		sch: ch,
		bch: make(chan []byte, 1000),
	}
}

func (e *EventFlow) Channel() chan []byte {
	safe.Go(func() {
		defer func() {
			if e.bch != nil {
				close(e.bch)
			}
		}()
		for {
			select {
			case event, ok := <-e.sch.Chan():
				if !ok {
					return
				}
				b, err := e.Serialize(event)
				if err != nil {
					log.Errorf("serialize event error: err: %v", err)
				}
				e.bch <- b
			}
		}
	})
	return e.bch
}

func (e *EventFlow) Serialize(event model.Event) ([]byte, error) {
	return json.Marshal(event)
}

func (e *EventFlow) dispatch(event model.Event) {
	e.sch.Send(event)
}

type dispatcher struct {
	mu  sync.RWMutex
	dst map[model.EngineGroup][]*EventFlow

	qSize   int
	timeout time.Duration
	ctx     context.Context
}

func newDispatcher(ctx context.Context) *dispatcher {
	return &dispatcher{
		dst: make(map[model.EngineGroup][]*EventFlow),

		qSize:   1000,
		timeout: time.Second * 3,
		ctx:     ctx,
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
