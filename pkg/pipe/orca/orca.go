package orca

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/yxxchange/pipefree/pkg/interfaces"
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

	*dispatcher
}

var _ interfaces.Orchestrator[*OrcaContext] = &Orchestrator{}

func GetOrchestrator(ctx context.Context) *Orchestrator {
	once.Do(func() {
		orca = &Orchestrator{
			ctx: &OrcaContext{
				ctx:       ctx,
				DAG:       make(map[string]*TopologyNode),
				PhaseRepo: make(map[model.Phase]*PhaseNodes),
			},
			dispatcher: newDispatcher(ctx),
		}
	})
	return orca
}

type OrcaContext struct {
	Pipe      model.Node                  `json:"node"`
	DAG       map[string]*TopologyNode    `json:"dag"`
	PhaseRepo map[model.Phase]*PhaseNodes `json:"phaseRepo"`

	ctx context.Context
}

type PhaseNodes struct {
	Phase model.Phase
	Map   map[string]*model.Node
}

func (m *Orchestrator) Serialize(ctx *OrcaContext) ([]byte, error) {
	// Serialize the node to JSON
	return json.Marshal(ctx.Pipe)
}

func (m *Orchestrator) Deserialize(b []byte) (*OrcaContext, error) {
	// Deserialize the JSON data to a Node
	var node model.Node
	err := json.Unmarshal(b, &node)
	if err != nil {
		return nil, err
	}
	m.ctx.Pipe = node
	err = m.init(m.ctx)
	if err != nil {
		return nil, err
	}
	return m.ctx, nil
}

func (m *Orchestrator) AddDispatcher(eg model.EngineGroup) *EventFlow {
	return m.dispatcher.Register(eg)
}

func (m *Orchestrator) init(ctx *OrcaContext) error {
	sorter, err := NewTopologySorter().ExtractGraph(ctx.Pipe.Graph).TopologySort()
	if err != nil {
		return err
	}
	ctx.DAG = sorter.Map
	readyNodes := sorter.GetZeroNode()
	m.ctx = ctx
	ctx.PhaseRepo[model.PhaseReady] = initPhaseRepo(readyNodes)
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

func (m *Orchestrator) Trigger(target string, action interfaces.Action) error {
	switch action {
	case interfaces.ActionStart:
		err := m.readyToRunning(target)
		if err != nil {
			return err
		}
		// TODO: dispatch to engine
		return nil

	default:
		return ErrorNotSupport
	}
}

func (m *Orchestrator) readyToRunning(target string) error {
	readyNodes := m.getPhaseNodes(model.PhaseReady)
	runningNodes := m.getPhaseNodes(model.PhaseRunning)
	node, has := readyNodes.search(target)
	if !has {
		return ErrorNodeNotReady
	}
	readyNodes.remove(target)
	node.Status.Phase = model.PhaseRunning
	runningNodes.add(target, node)
	return nil
}

func (m *Orchestrator) getPhaseNodes(phase model.Phase) *PhaseNodes {
	if m.ctx.PhaseRepo == nil {
		m.ctx.PhaseRepo = make(map[model.Phase]*PhaseNodes)
	}
	return m.ctx.PhaseRepo[phase]
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
