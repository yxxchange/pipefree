package orca

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
)

type TopologyNodes []*TopologyNode

var (
	ErrorCycle         = errors.New("graph has a cycle")
	ErrorEmpty         = errors.New("graph is empty")
	ErrorInvalidEdge   = errors.New("graph edge is invalid")
	ErrorOriginMustOne = errors.New("the number of origin node must be one")
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

type GraphParser struct {
	Map    map[string]*TopologyNode
	List   TopologyNodes
	backup TopologyNodes

	Err error

	parsed bool
}

func NewGraphParser() *GraphParser {
	return &GraphParser{}
}

// Parse before calling this method, you must ensure that pipe has been prepared
func (t *GraphParser) Parse(pipe model.PipeFlow) *GraphParser {
	t.parsed = true
	t.Map = make(map[string]*TopologyNode)
	t.List = make(TopologyNodes, 0)
	t.parse(pipe)
	return t
}

func (t *GraphParser) FindTheOrigin() (model.Node, error) {
	if !t.parsed {
		return model.Node{}, fmt.Errorf("graph not parsed")
	}
	for _, node := range t.backup {
		if node.InDegree == 0 {
			return node.Node, nil
		}
	}
	return model.Node{}, ErrorOriginMustOne
}

func (t *GraphParser) parse(pipe model.PipeFlow) {
	for _, node := range pipe.Nodes {
		topoNode := &TopologyNode{
			Node:     node,
			InDegree: 0,
		}
		t.List = append(t.List, topoNode)
		t.Map[node.MetaData.Name] = topoNode
	}
	for _, edge := range pipe.Graph.Edges {
		fromNode := t.Map[edge.From]
		toNode := t.Map[edge.To]
		if fromNode == nil || toNode == nil {
			t.Err = ErrorInvalidEdge
			return
		}
		toNode.InDegree++
		fromNode.AddToNode(&toNode.MetaData)
		toNode.AddFromNode(&fromNode.MetaData)
	}
	if len(t.List) == 0 {
		t.Err = ErrorEmpty
	}
	cnt := 0
	for _, node := range t.backup {
		if node.InDegree == 0 {
			cnt++
		}
	}
	if cnt > 1 {
		t.Err = ErrorOriginMustOne
		return
	}
	t.backup = make(TopologyNodes, len(t.List))
	copy(t.backup, t.List)
}

func (t *GraphParser) IsValid() error {
	if !t.parsed {
		return fmt.Errorf("graph not parsed")
	}
	return t.sort()
}

func (t *GraphParser) sort() error {
	if t.Err != nil {
		return t.Err
	}
	heap.Init(&t.List)
	if len(t.List) == 0 {
		return ErrorEmpty
	}
	for len(t.List) > 0 {
		node := heap.Pop(&t.List).(*TopologyNode)
		if node.InDegree != 0 {
			return ErrorCycle
		}
		for _, nextMeta := range node.Node.MetaData.To {
			nextNode := t.Map[nextMeta.Name]
			nextNode.InDegree--
		}
		heap.Init(&t.List)
	}
	return nil
}

type TopologyNode struct {
	model.Node `json:"node"`
	InDegree   int `json:"inDegree"`
}
