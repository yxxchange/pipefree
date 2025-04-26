package orca

import (
	"container/heap"
	"errors"
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

type TopologySorter struct {
	Map  map[string]*TopologyNode
	List TopologyNodes

	Backup TopologyNodes // static data, used to backup the original data
}

func NewTopologySorter() *TopologySorter {
	return &TopologySorter{}
}

func (t *TopologySorter) GetZeroNode() TopologyNodes {
	var res TopologyNodes
	heap.Init(&t.List)

	for len(t.List) > 0 {
		node := heap.Pop(&t.List).(*TopologyNode)
		if node.InDegree == 0 {
			res = append(res, node)
			continue
		}
		heap.Push(&t.List, node)
		break
	}

	for _, node := range res {
		successor := t.Map[node.Node.Name].Node.GetSuccessor()
		for _, s := range successor {
			t.Map[s.Name].InDegree--
		}
	}

	return res
}

func (t *TopologySorter) TopologySort() (*TopologySorter, error) {
	zeroInDegree, mapCopy := t.deepCopy()
	if len(zeroInDegree) == 0 {
		if len(t.Map) == 0 {
			return nil, ErrorEmpty
		}
		return nil, ErrorCycle
	}
	heap.Init(&zeroInDegree)
	for size := len(zeroInDegree); size > 0; size = len(zeroInDegree) {
		for i := 0; i < size; i++ {
			node := heap.Pop(&zeroInDegree).(*TopologyNode)
			if node.InDegree > 0 {
				return nil, ErrorCycle
			}
			t.store(t.Map[node.Node.Name])
			for _, successor := range node.Node.GetSuccessor() {
				successorNode := mapCopy[successor.Name]
				successorNode.InDegree--
				if successorNode.InDegree == 0 {
					heap.Push(&zeroInDegree, successorNode)
				}
			}
		}
	}

	if len(t.List) != len(t.Map) {
		return nil, ErrorCycle
	}

	return t, nil
}

func (t *TopologySorter) ExtractGraph(graph model.Graph) *TopologySorter {
	nodesMap := make(map[string]*TopologyNode)
	for _, node := range graph.Nodes {
		nodesMap[node.Name] = &TopologyNode{
			Node:     &node,
			InDegree: 0,
		}
	}
	for _, edge := range graph.Edges {
		from, ok1 := nodesMap[edge.From]
		to, ok2 := nodesMap[edge.To]
		if ok1 && ok2 {
			from.Node.AddSuccessor(&to.Node.MetaData)
			to.Node.AddPredecessor(&from.Node.MetaData)
		}
		if ok2 {
			to.InDegree++
		}
	}
	t.Map = nodesMap
	return t
}

func (t *TopologySorter) deepCopy() (TopologyNodes, map[string]*TopologyNode) {
	zeroInDegree := make(TopologyNodes, 0, len(t.Map))
	copied := make(map[string]*TopologyNode, len(t.Map))
	for _, node := range t.Map {
		backup := node.DeepCopy()
		if backup.InDegree == 0 {
			zeroInDegree = append(zeroInDegree, backup)
		}
		copied[node.Node.Name] = backup
	}
	return zeroInDegree, copied
}

func (t *TopologySorter) store(node *TopologyNode) {
	t.List = append(t.List, node)
	t.Backup = append(t.Backup, node.DeepCopy())
}

type TopologyNode struct {
	// Node is the node of the orca
	Node     *model.Node `json:"node"`
	InDegree int         `json:"inDegree"`
}

func (tn TopologyNode) DeepCopy() *TopologyNode {
	return &TopologyNode{
		Node:     tn.Node,
		InDegree: tn.InDegree,
	}
}
