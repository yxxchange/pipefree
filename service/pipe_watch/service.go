package pipe_watch

import (
	"context"
	"github.com/yxxchange/pipefree/helper/safe"
)

type Service struct {
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) Watch(operator OperatorCtx) *EventChannel {
	safe.Go(func() {
		GetWatchServer().ListAndWatch(operator)
	})
	return operator.EventChannel
}

func (s *Service) UnWatch(operator OperatorCtx) {
	GetWatchServer().RemoveOperator(operator)
}
