package method

import (
	"context"
	"encoding/json"
	"errors"
	"pipefree/pkg/interfaces"
	"pipefree/pkg/pipe/data"
	"sync"
)

const (
	OrcaCtxKey = "orca"
)

var (
	ErrorNotSupport   = errors.New("not support")
	ErrorNodeNotReady = errors.New("node not in ready")
)

type Orchestrator struct {
	ctx *OrcaCtx
	mu  sync.Mutex
}

var _ interfaces.Orchestrator[*OrcaCtx] = &Orchestrator{}

func NewOrchestrator(ctx context.Context) *Orchestrator {
	return &Orchestrator{
		ctx: &OrcaCtx{
			ctx:       ctx,
			DAG:       make(map[string]*TopologyNode),
			PhaseRepo: make(map[data.Phase]*PhaseNodes),
		},
	}
}

type OrcaCtx struct {
	Pipe      data.Node                  `json:"node"`
	DAG       map[string]*TopologyNode   `json:"dag"`
	PhaseRepo map[data.Phase]*PhaseNodes `json:"phaseRepo"`

	ctx context.Context
}

type PhaseNodes struct {
	Phase data.Phase
	Map   map[string]*data.Node
}

func (m *Orchestrator) Serialize(ctx *OrcaCtx) ([]byte, error) {
	// Serialize the node to JSON
	return json.Marshal(ctx.Pipe)
}

func (m *Orchestrator) Deserialize(b []byte) (*OrcaCtx, error) {
	// Deserialize the JSON data to a Node
	var node data.Node
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

func (m *Orchestrator) init(ctx *OrcaCtx) error {
	sorter, err := NewTopologySorter().ExtractGraph(ctx.Pipe.Graph).TopologySort()
	if err != nil {
		return err
	}
	ctx.DAG = sorter.Map
	readyNodes := sorter.GetZeroNode()
	m.ctx = ctx
	ctx.PhaseRepo[data.PhaseReady] = initPhaseRepo(readyNodes)
	return nil
}

func initPhaseRepo(readyNodes TopologyNodes) *PhaseNodes {
	phaseNodes := PhaseNodes{
		Phase: data.PhaseReady,
		Map:   make(map[string]*data.Node),
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
	readyNodes := m.getPhaseNodes(data.PhaseReady)
	runningNodes := m.getPhaseNodes(data.PhaseRunning)
	node, has := readyNodes.search(target)
	if !has {
		return ErrorNodeNotReady
	}
	readyNodes.remove(target)
	node.Status.Phase = data.PhaseRunning
	runningNodes.add(target, node)
	return nil
}

func (m *Orchestrator) getPhaseNodes(phase data.Phase) *PhaseNodes {
	if m.ctx.PhaseRepo == nil {
		m.ctx.PhaseRepo = make(map[data.Phase]*PhaseNodes)
	}
	return m.ctx.PhaseRepo[phase]
}

func (pn *PhaseNodes) search(name string) (*data.Node, bool) {
	if len(pn.Map) == 0 {
		return nil, false
	}
	node, has := pn.Map[name]
	if !has {
		return nil, false
	}
	return node, true
}

func (pn *PhaseNodes) add(name string, node *data.Node) {
	if pn.Map == nil {
		pn.Map = make(map[string]*data.Node)
	}
	pn.Map[name] = node
}

func (pn *PhaseNodes) remove(name string) {
	if pn.Map == nil {
		pn.Map = make(map[string]*data.Node)
		return
	}
	delete(pn.Map, name)
}
