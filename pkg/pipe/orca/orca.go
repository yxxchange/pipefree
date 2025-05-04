package orca

import (
	"context"
	"errors"
	"fmt"
	"github.com/yxxchange/pipefree/pkg/infra/nebula"
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

func (o *Orchestrator) Schedule(pipe model.Node) error {
	return nil
}

func (o *Orchestrator) Create(raw string) error {
	vertexes, err := NewGraphBuilder().ProcessYaml(raw).ProcessGraph().Build()
	if err != nil {
		return fmt.Errorf("create pipe error: %v", err)
	}
	if len(vertexes) == 0 {
		return fmt.Errorf("pipe is empty")
	}
	// todo: sql
	nebula.UseSpace(vertexes[0].MetaData.Space).Execute("")
	return nil
}
