package orca

import (
	"context"
	"errors"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"sync"
)

var orca *Orchestrator
var once sync.Once

var (
	ErrorNotSupport   = errors.New("not support")
	ErrorNodeNotReady = errors.New("node not in ready")
)

type Orchestrator struct {
	*watcher
}

func GetOrchestrator(ctx context.Context) *Orchestrator {
	once.Do(func() {
		orca = &Orchestrator{
			watcher: newWatcher(ctx),
		}
	})
	return orca
}

func (o *Orchestrator) Schedule(pipe model.PipeFlow) error {
	return nil
}
