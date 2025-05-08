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
	View     model.Node
	Map      map[string]*TopologyNode
	List     TopologyNodes
	Vertexes []*model.NodeInfo

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

func (t *GraphBuilder) Build() ([]*model.NodeInfo, error) {
	if t.Err != nil {
		return nil, t.Err
	}
	var vertexes []*model.NodeInfo
	for _, topo := range t.List {
		vertexes = append(vertexes, topo.Node)
	}
	return vertexes, nil
}

func (t *GraphBuilder) ProcessYaml(raw string) *GraphBuilder {
	if t.Err != nil {
		return t
	}
	var nodeView model.Node
	if err := serialize.YamlDeserialize([]byte(raw), &nodeView); err != nil {
		t.Err = fmt.Errorf("deserialize pipe error: %v", err)
		return t
	}
	t.View = nodeView
	t.validate(&t.View)
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

func (t *GraphBuilder) validate(node *model.Node) *GraphBuilder {
	if t.Err != nil {
		return t
	}
	if node.MetaData.Name == "" {
		t.Err = fmt.Errorf("node name is empty")
		return t
	}
	if _, ok := t.nameCollector[node.MetaData.Name]; ok {
		t.Err = fmt.Errorf("node name %s is duplicated", node.MetaData.Name)
		return t
	}
	if node.ApiVersion == "" {
		t.Err = fmt.Errorf("node apiVersion is empty")
		return t
	}
	if node.Kind == model.NodeKindScalar && node.Graph != nil {
		t.Err = fmt.Errorf("node of scalar kind should not contain graph defined")
		return t
	}
	if node.Kind == model.NodeKindCompound && node.Graph == nil {
		t.Err = fmt.Errorf("node of compound kind should contain graph defined")
		return t
	}
	if node.Graph != nil {
		for i := 0; i < len(node.Graph.Vertexes); i++ {
			t.validate(&node.Graph.Vertexes[i])
			if t.Err != nil {
				return t
			}
		}
	}
	return t
}

func (t *GraphBuilder) ProcessGraph() *GraphBuilder {
	if t.Err != nil {
		return t
	}
	graph := t.View.Graph
	t.Map = make(map[string]*TopologyNode)
	// add the head node
	t.Map[t.View.MetaData.Name] = &TopologyNode{
		Node:     &t.View.NodeInfo,
		InDegree: 0,
	}
	for i := 0; i < len(graph.Vertexes); i++ {
		node := &TopologyNode{
			Node:     &graph.Vertexes[i].NodeInfo,
			InDegree: 0,
		}
		t.Map[graph.Vertexes[i].MetaData.Name] = node
	}
	for i := 0; i < len(graph.Edges); i++ {
		if _, ok := t.Map[graph.Edges[i].From]; !ok {
			t.Err = fmt.Errorf("node %s not found", graph.Edges[i].From)
			return t
		}
		if _, ok := t.Map[graph.Edges[i].To]; !ok {
			t.Err = fmt.Errorf("node %s not found", graph.Edges[i].To)
			return t
		}
		t.Map[graph.Edges[i].To].InDegree++
		t.Map[graph.Edges[i].From].Node.MetaData.AddTo(&t.Map[graph.Edges[i].To].Node.MetaData)
		t.Map[graph.Edges[i].To].Node.MetaData.AddFrom(&t.Map[graph.Edges[i].From].Node.MetaData)
		t.Map[graph.Edges[i].To].Node.MetaData.AddAncestor(&t.View.MetaData)
	}
	if len(t.Map) == 0 {
		t.Err = ErrorEmpty
		return t
	}
	t.List = make(TopologyNodes, 0, len(t.Map))
	for _, node := range t.Map {
		t.List = append(t.List, node)
	}
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
	Node     *model.NodeInfo `json:"node"`
	InDegree int             `json:"inDegree"`
}
