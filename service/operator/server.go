package operator

import (
	"context"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/infra/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var serverInstance *WatchServer
var once sync.Once

func GetWatchServer() *WatchServer {
	once.Do(func() {
		serverInstance = NewWatchServer()
		go serverInstance.Monitor()
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

func (s *WatchServer) Monitor() {
	ticker := time.NewTimer(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.lock.Lock()
			log.Infof("当前监控的前缀数量: %d", len(s.prefixMap))
			for prefix, handler := range s.prefixMap {
				if handler == nil {
					log.Warnf("前缀 %s 的处理器为 nil，可能已被删除", prefix)
					delete(s.prefixMap, prefix) // 删除无效的处理器
					continue
				}
				log.Infof("前缀: %s, 监控通道数量: %d", prefix, len(handler.channels))
			}
		case <-s.ctx.Done():
			return // 退出监控
		}
	}
}

type WatchHandler struct {
	lock     sync.RWMutex
	ctx      context.Context
	prefix   string
	channels []*EventChannel
}

func NewWatchHandler(ctx context.Context, prefix string) *WatchHandler {
	return &WatchHandler{
		ctx:      ctx,
		prefix:   prefix,
		channels: make([]*EventChannel, 100),
	}
}

func (s *WatchServer) Register(prefix string, ch *EventChannel) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.prefixMap == nil {
		s.prefixMap = make(map[string]*WatchHandler)
	}
	if handler, exists := s.prefixMap[prefix]; exists {
		handler.Register(ch)
	} else {
		s.prefixMap[prefix] = NewWatchHandler(s.ctx, prefix)
		s.prefixMap[prefix].Register(ch)
	}
}

func (h *WatchHandler) Register(ch *EventChannel) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.channels == nil {
		h.channels = make([]*EventChannel, 0, 100) // 初始化通道切片
		go h.Watch()                               // 启动监听
	}
	h.channels = append(h.channels, ch)
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
		}
		return
	}

	for _, e := range result.Events {
		event := Convert(e)
		for _, ch := range h.channels {
			ch.ch <- event
		}
	}
}
