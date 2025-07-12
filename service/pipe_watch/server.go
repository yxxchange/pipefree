package pipe_watch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/helper/safe"
	"github.com/yxxchange/pipefree/infra/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

var serverInstance *WatchServer
var once sync.Once

func GetWatchServer() *WatchServer {
	once.Do(func() {
		serverInstance = NewWatchServer()
	})
	return serverInstance
}

type WatchServer struct {
	lock      sync.Mutex
	ctx       context.Context
	streamAgg map[string]*EventStream
}

func NewWatchServer() *WatchServer {
	return &WatchServer{
		ctx:       context.Background(),
		streamAgg: make(map[string]*EventStream),
	}
}

type EventStream struct {
	lock   sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	watchStarted bool // 是否已启动监听
	revSince     int64
	streamId     string
	channels     map[string]*EventChannel
}

func NewEventStream(ctx context.Context, prefix string) *EventStream {
	son, cancel := context.WithCancel(ctx)
	return &EventStream{
		cancel:   cancel,
		ctx:      son,
		streamId: prefix,
	}
}

func (s *WatchServer) ListAndWatch(operator OperatorCtx) {
	stream := s.GetEventStream(operator.StreamID)
	err := stream.ListAndWatch(operator.UUID, operator.EventChannel)
	if err != nil { // 处理错误
		errMsg := fmt.Errorf("failed to list and watch streamId %s: %v", stream.streamId, err)
		operator.EventChannel.SendErr(errMsg)
		operator.EventChannel.Close()
		return
	}
}

func (s *WatchServer) GetEventStream(prefix string) *EventStream {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.streamAgg == nil {
		s.streamAgg = make(map[string]*EventStream)
	}
	if stream, exists := s.streamAgg[prefix]; exists {
		return stream
	}
	stream := NewEventStream(s.ctx, prefix)
	s.streamAgg[prefix] = stream
	return stream
}

func (s *WatchServer) RemoveOperator(operator OperatorCtx) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if stream, exists := s.streamAgg[operator.StreamID]; exists {
		stream.RemoveChannel(operator.UUID)
		if len(stream.channels) == 0 {
			log.Infof("No more channels for streamId %s, stopping watch", operator.StreamID)
			stream.cancel()
			delete(s.streamAgg, operator.StreamID) // 删除无效的事件流
		}
	}
}

func (h *EventStream) RemoveChannel(uuid string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.channels != nil {
		if _, exists := h.channels[uuid]; exists {
			delete(h.channels, uuid) // 删除通道
		}
	}
}

func (h *EventStream) ListAndWatch(uuid string, ch *EventChannel) error {
	rev, err := h.List(ch)
	if err != nil {
		return fmt.Errorf("failed to list streamId %s: %v", h.streamId, err)
	}
	h.AddChannel(uuid, ch)
	if !h.watchStarted {
		h.watchStarted = true // 标记监听已启动
		h.revSince = rev
		safe.Go(h.Watch)
	}
	return nil
}

func (h *EventStream) AddChannel(uuid string, ch *EventChannel) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.channels == nil {
		h.channels = make(map[string]*EventChannel) // 初始化通道切片
	}
	h.channels[uuid] = ch
}

func (h *EventStream) Remove(uuid string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.channels, uuid)
}

func (h *EventStream) List(ch *EventChannel) (int64, error) {
	resp, err := etcd.GetWithPrefix(h.ctx, h.streamId)
	if err != nil {
		log.Errorf("failed to get streamId %s: %v", h.streamId, err)
		return 0, err
	}
	h.HandleList(resp, ch)
	return resp.Header.Revision, nil
}

func (h *EventStream) Watch() {
	etcd.Watch(h.ctx, h.streamId, h.revSince, h.HandleWatch)
}

func (h *EventStream) HandleList(result *clientv3.GetResponse, ch *EventChannel) {
	for _, kv := range result.Kvs {
		event := Convert(&clientv3.Event{
			Type: clientv3.EventTypePut,
			Kv:   kv,
		})
		b, err := json.Marshal(event)
		if err != nil {
			log.Errorf("failed to marshal event: %v", err)
			continue // 如果序列化失败，跳过当前事件
		}
		ch.ch <- b // 发送事件到通道
	}
}

func (h *EventStream) HandleWatch(result *clientv3.WatchResponse, closed bool) {
	if closed {
		log.Infof("watch closed for streamId %s", h.streamId)
		for _, ch := range h.channels {
			ch.done <- struct{}{}
			ch.Close() // 关闭通道
		}
		return
	}

	for _, e := range result.Events {
		event := Convert(e)
		b, err := json.Marshal(event)
		if err != nil {
			log.Errorf("failed to marshal event: %v", err)
			continue // 如果序列化失败，跳过当前事件
		}
		for _, ch := range h.channels {
			ch.ch <- b
		}
	}
}
