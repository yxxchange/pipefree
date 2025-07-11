package operator

import (
	"context"
	"encoding/json"
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
	prefixMap map[string]*WatchHandler
}

func NewWatchServer() *WatchServer {
	return &WatchServer{
		ctx:       context.Background(),
		prefixMap: make(map[string]*WatchHandler),
	}
}

type WatchHandler struct {
	lock     sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	prefix   string
	channels map[string]*EventChannel
}

func NewWatchHandler(ctx context.Context, prefix string) *WatchHandler {
	son, cancel := context.WithCancel(ctx)
	return &WatchHandler{
		cancel: cancel,
		ctx:    son,
		prefix: prefix,
	}
}

func (s *WatchServer) Register(prefix, uuid string, ch *EventChannel) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.prefixMap == nil {
		s.prefixMap = make(map[string]*WatchHandler)
	}
	if handler, exists := s.prefixMap[prefix]; exists {
		handler.Register(uuid, ch)
	} else {
		s.prefixMap[prefix] = NewWatchHandler(s.ctx, prefix)
		s.prefixMap[prefix].Register(uuid, ch)
	}
}

func (s *WatchServer) UnRegister(prefix, uuid string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if handler, exists := s.prefixMap[prefix]; exists {
		handler.UnRegister(uuid)
		if len(handler.channels) == 0 {
			log.Infof("No more channels for prefix %s, stopping watch", prefix)
			delete(s.prefixMap, prefix) // 删除无效的处理器
		}
	}
}

func (h *WatchHandler) Register(uuid string, ch *EventChannel) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.channels == nil {
		h.channels = make(map[string]*EventChannel) // 初始化通道切片
		safe.Go(h.Watch)                            // 启动监听
	}
	h.channels[uuid] = ch
}

func (h *WatchHandler) UnRegister(uuid string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.channels != nil {
		if _, exists := h.channels[uuid]; exists {
			delete(h.channels, uuid) // 删除通道
		}
		if len(h.channels) == 0 {
			log.Infof("No more channels for prefix %s, stopping watch", h.prefix)
			h.cancel()
		}
	}
}

func (h *WatchHandler) Watch() {
	etcd.Watch(h.ctx, h.prefix, h.Handle)
}

func (h *WatchHandler) Handle(result *clientv3.WatchResponse, closed bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if closed {
		log.Infof("watch closed for prefix %s", h.prefix)
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
