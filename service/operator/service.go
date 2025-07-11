package operator

import (
	"context"
)

type Service struct {
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) Watch(keyPrefix, uuid string) *EventChannel {
	ch := &EventChannel{
		ch:   make(chan []byte, 100), // 设置缓冲区大小为100
		done: make(chan struct{}),
	}
	GetWatchServer().Register(keyPrefix, uuid, ch)
	return ch
}

func (s *Service) UnWatch(keyPrefix, uuid string) {
	GetWatchServer().UnRegister(keyPrefix, uuid)
}
