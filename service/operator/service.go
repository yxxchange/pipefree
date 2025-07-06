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

func (s *Service) Watch(keyPrefix string) *EventChannel {
	ch := &EventChannel{
		ch: make(chan Event, 100), // 设置缓冲区大小为100
	}
	GetWatchServer().Register(keyPrefix, ch)
	return ch
}
