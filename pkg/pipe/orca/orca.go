package orca

import (
	"context"
	"errors"
	"fmt"
	"github.com/yxxchange/pipefree/helper/log"
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

type OrcaContext struct {
	Pipe model.Node               `json:"pipe"` // pipe flow from yaml
	DAG  map[string]*TopologyNode `json:"dag"`
	//PhaseRepo map[model.Phase]*PhaseNodes `json:"phaseRepo"`

	ctx context.Context
}

type PhaseNodes struct {
	Phase model.Phase
	Map   map[string]*model.Node
}

func (o *Orchestrator) Schedule(pipe model.Node) error {
	ctx := &OrcaContext{
		ctx:  context.TODO(),
		Pipe: pipe,
	}
	err := o.setup(ctx)
	if err != nil {
		log.Errorf("setup error when schedule: %v", err)
		return err
	}
	return nil
}

func (o *Orchestrator) setup(ctx *OrcaContext) error {
	// TODO: 需要重构
	sorter, err := NewTopologySorter().ExtractGraph(ctx.Pipe.Graph).TopologySort()
	if err != nil {
		return err
	}
	ctx.DAG = sorter.Map
	headNode := sorter.GetZeroNode()
	if len(headNode) != 1 {
		return fmt.Errorf("unexpected result when setup pipe, expected 1, but %d", len(headNode))
	}

	return nil
}

func (pn *PhaseNodes) search(name string) (*model.Node, bool) {
	if len(pn.Map) == 0 {
		return nil, false
	}
	node, has := pn.Map[name]
	if !has {
		return nil, false
	}
	return node, true
}

func (pn *PhaseNodes) add(name string, node *model.Node) {
	if pn.Map == nil {
		pn.Map = make(map[string]*model.Node)
	}
	pn.Map[name] = node
}

func (pn *PhaseNodes) remove(name string) {
	if pn.Map == nil {
		pn.Map = make(map[string]*model.Node)
		return
	}
	delete(pn.Map, name)
}
