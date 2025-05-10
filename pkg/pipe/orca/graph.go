package orca

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/yxxchange/pipefree/helper/serialize"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
)

type TopologyNodes []*TopologyNode

var (
	ErrorCycle = errors.New("graph has a cycle")
	ErrorEmpty = errors.New("graph is empty")
)

func (t TopologyNodes) Len() int {
	return len(t)
}

func (t TopologyNodes) Less(i, j int) bool {
	return t[i].InDegree < t[j].InDegree
}

func (t TopologyNodes) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t *TopologyNodes) Push(x any) {
	*t = append(*t, x.(*TopologyNode))
}

func (t *TopologyNodes) Pop() any {
	x := (*t)[len(*t)-1]
	*t = (*t)[:len(*t)-1]
	return x
}

type GraphBuilder struct {
	Node     model.PipeFlow
	Map      map[string]*TopologyNode
	List     TopologyNodes
	Vertexes []*model.Node

	nameCollector  map[string]struct{}
	spaceCollector map[string]struct{}
	tagCollector   map[string]struct{}

	Err error
}

func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		nameCollector:  make(map[string]struct{}),
		spaceCollector: make(map[string]struct{}),
		tagCollector:   make(map[string]struct{}),
	}
}

func (t *GraphBuilder) Build() ([]*model.Node, error) {
	if t.Err != nil {
		return nil, t.Err
	}
	var vertexes []*model.Node
	for _, topo := range t.List {
		vertexes = append(vertexes, topo.Node)
	}
	return vertexes, nil
}

func (t *GraphBuilder) ProcessPipeCfg(pipe model.PipeConfig) *GraphBuilder {
	if t.Err != nil {
		return t
	}
	t.Node = pipe.PipeFlow
	if len(t.spaceCollector) > 1 {
		t.Err = fmt.Errorf("only one space is allowed, but got %d", len(t.spaceCollector))
		return t
	}
	if len(t.tagCollector) > 1 {
		t.Err = fmt.Errorf("only one tag is allowed, but got %d", len(t.tagCollector))
		return t
	}
	return t
}

func (t *GraphBuilder) ProcessYaml(raw string) *GraphBuilder {
	if t.Err != nil {
		return t
	}
	var nodeView model.PipeFlow
	if err := serialize.YamlDeserialize([]byte(raw), &nodeView); err != nil {
		t.Err = fmt.Errorf("deserialize pipe error: %v", err)
		return t
	}
	t.Node = nodeView
	if len(t.spaceCollector) > 1 {
		t.Err = fmt.Errorf("only one space is allowed, but got %d", len(t.spaceCollector))
		return t
	}
	if len(t.tagCollector) > 1 {
		t.Err = fmt.Errorf("only one tag is allowed, but got %d", len(t.tagCollector))
		return t
	}
	return t
}

func (t *GraphBuilder) ProcessGraph() *GraphBuilder {
	if t.Err != nil {
		return t
	}
	// todo
	return t.sort()
}

func (t *GraphBuilder) sort() *GraphBuilder {
	if t.Err != nil {
		return t
	}
	heap.Init(&t.List)
	for len(t.List) > 0 {
		node := heap.Pop(&t.List).(*TopologyNode)
		if node.InDegree != 0 {
			t.Err = ErrorCycle
			return t
		}
		for _, nextMeta := range node.Node.MetaData.To {
			nextNode := t.Map[nextMeta.Name]
			nextNode.InDegree--
			if nextNode.InDegree == 0 {
				heap.Push(&t.List, nextNode)
			}
		}
	}
	return t
}

type TopologyNode struct {
	// Node is the node of the orca
	Node     *model.Node `json:"node"`
	InDegree int         `json:"inDegree"`
}
