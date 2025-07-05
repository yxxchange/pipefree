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

func ServerInstance() *WatchServer {
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
				log.Infof("前缀: %s, 监控通道数量: %d", prefix, len(handler.chMap))
				operators := make([]string, 0, len(handler.chMap))
				for operatorName := range handler.chMap {
					operators = append(operators, operatorName)
				}
				log.Infof("前缀: %s, 监控操作员: %v", prefix, operators)
			}
		case <-s.ctx.Done():
			return // 退出监控
		}
	}
}

type WatchHandler struct {
	lock   sync.RWMutex
	ctx    context.Context
	prefix string
	chMap  map[string]*EventChannel
}

func NewWatchHandler(ctx context.Context, prefix string) *WatchHandler {
	return &WatchHandler{
		ctx:    ctx,
		prefix: prefix,
		chMap:  make(map[string]*EventChannel),
	}
}

func (s *WatchServer) Register(prefix, operatorName string, ch *EventChannel) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.prefixMap == nil {
		s.prefixMap = make(map[string]*WatchHandler)
	}
	if handler, exists := s.prefixMap[prefix]; exists {
		handler.Register(operatorName, ch)
	} else {
		s.prefixMap[prefix] = NewWatchHandler(s.ctx, prefix)
		s.prefixMap[prefix].Register(operatorName, ch)
	}
}

func (h *WatchHandler) Register(uuid string, ch *EventChannel) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.chMap == nil {
		h.chMap = make(map[string]*EventChannel)
		go h.Watch() // 启动监听
	}
	h.chMap[uuid] = ch
}

func (h *WatchHandler) Watch() {
	etcd.Watch(h.ctx, h.prefix, h.Handle)
}

func (h *WatchHandler) Handle(result *clientv3.WatchResponse) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for _, e := range result.Events {
		event := Convert(e)
		for _, ch := range h.chMap {
			ch.ch <- event
		}
	}
}
