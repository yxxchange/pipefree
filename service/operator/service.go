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

func (s *Service) Watch(keyPrefix, operatorName string) *EventChannel {
	ch := &EventChannel{
		ch: make(chan Event, 100), // 设置缓冲区大小为100
	}
	ServerInstance().Register(keyPrefix, operatorName, ch)
	return ch
}
