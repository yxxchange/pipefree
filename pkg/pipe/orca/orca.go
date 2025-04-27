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
	ctx *OrcaContext
	mu  sync.Mutex

	*watcher
}

func GetOrchestrator(ctx context.Context) *Orchestrator {
	once.Do(func() {
		orca = &Orchestrator{
			ctx: &OrcaContext{
				ctx:       ctx,
				DAG:       make(map[string]*TopologyNode),
				PhaseRepo: make(map[model.Phase]*PhaseNodes),
			},
			watcher: newWatcher(ctx),
		}
	})
	return orca
}

type OrcaContext struct {
	Pipe      model.Node                  `json:"pipe"` // pipe flow from yaml
	DAG       map[string]*TopologyNode    `json:"dag"`
	PhaseRepo map[model.Phase]*PhaseNodes `json:"phaseRepo"`

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
	return o.init(ctx)
}

func (o *Orchestrator) init(ctx *OrcaContext) error {
	sorter, err := NewTopologySorter().ExtractGraph(ctx.Pipe.Graph).TopologySort()
	if err != nil {
		return err
	}
	ctx.DAG = sorter.Map
	readyNodes := sorter.GetZeroNode()
	ctx.PhaseRepo[model.PhaseReady] = initPhaseRepo(readyNodes)
	o.ctx = ctx
	return nil
}

func initPhaseRepo(readyNodes TopologyNodes) *PhaseNodes {
	phaseNodes := PhaseNodes{
		Phase: model.PhaseReady,
		Map:   make(map[string]*model.Node),
	}
	for _, node := range readyNodes {
		phaseNodes.Map[node.Node.Name] = node.Node
	}
	return &phaseNodes
}

func (o *Orchestrator) readyToRunning(target string) error {
	readyNodes := o.getPhaseNodes(model.PhaseReady)
	runningNodes := o.getPhaseNodes(model.PhaseRunning)
	node, has := readyNodes.search(target)
	if !has {
		return ErrorNodeNotReady
	}
	readyNodes.remove(target)
	node.Status.Phase = model.PhaseRunning
	runningNodes.add(target, node)
	return nil
}

func (o *Orchestrator) getPhaseNodes(phase model.Phase) *PhaseNodes {
	if o.ctx.PhaseRepo == nil {
		o.ctx.PhaseRepo = make(map[model.Phase]*PhaseNodes)
	}
	return o.ctx.PhaseRepo[phase]
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
