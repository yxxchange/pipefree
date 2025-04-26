package orca

import (
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/helper/safe"
	"github.com/yxxchange/pipefree/pkg/infra/etcd"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EventFlow struct {
	bytesCh *safe.Channel[[]byte]
	eventCh *safe.Channel[*clientv3.WatchResponse]
}

func NewEventFlow(ch *safe.Channel[*clientv3.WatchResponse]) *EventFlow {
	return &EventFlow{
		eventCh: ch,
		bytesCh: safe.NewSafeChannel[[]byte](viper.GetInt("orca.dispatcher.queueSize")),
	}
}

func (e *EventFlow) Channel() *safe.Channel[[]byte] {
	safe.Go(func() {
		defer e.bytesCh.Close()
		e.eventCh.Range(e.transform)
	})
	return e.bytesCh
}

func (e *EventFlow) Serialize(resp *clientv3.WatchResponse) ([]byte, error) {
	return json.Marshal(resp)
}

func (e *EventFlow) transform(resp *clientv3.WatchResponse) (interrupted bool) {
	if resp == nil {
		return
	}

	for _, etcdEvent := range resp.Events {
		var event model.Event
		switch etcdEvent.Type {
		case mvccpb.DELETE:
			event.EventType = model.EventTypeDelete
		default:
			if etcdEvent.Kv.CreateRevision == etcdEvent.Kv.ModRevision {
				event.EventType = model.EventTypeCreate
			} else {
				event.EventType = model.EventTypeUpdate
			}
		}
		event.Data = etcdEvent.Kv.Value
		b, err := e.Serialize(resp)
		if err != nil {
			log.Errorf("serialize event error: err: %v", err)
			continue
		}
		e.bytesCh.Send(b)
	}

	return
}

type dispatcher struct {
	ctx context.Context
}

func newDispatcher(ctx context.Context) *dispatcher {
	return &dispatcher{
		ctx: ctx,
	}
}

func (d *dispatcher) Register(idf model.NodeIdentifier) *EventFlow {
	ch := safe.NewSafeChannel[*clientv3.WatchResponse](viper.GetInt("orca.dispatcher.queueSize"))
	ef := NewEventFlow(ch)
	safe.Go(func() {
		for resp := range etcd.Watch(context.Background(), idf.Identifier()) {
			ch.Send(&resp)
		}
		log.Info("etch watch closed")
	})
	return ef
}

func (d *dispatcher) Dispatch(ctx context.Context, key, value string) error {
	return etcd.Put(ctx, key, value)
}
