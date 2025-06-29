package graph

import (
	"container/heap"
	"fmt"
	"github.com/yxxchange/pipefree/infra/dal/model"
)

type Vertexes []*Vertex

type Vertex struct {
	Name     string              `json:"-" yaml:"-"` // 节点名称
	InDegree int                 `json:"-" yaml:"-"` // 入度
	Next     map[string]struct{} `json:"-" yaml:"-"` // 下一个节点名称列表
}

var _ heap.Interface = (*Vertexes)(nil)

func (vs Vertexes) Len() int {
	return len(vs)
}

func (vs Vertexes) Less(i, j int) bool {
	return vs[i].InDegree < vs[j].InDegree
}

func (vs Vertexes) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs *Vertexes) Push(x interface{}) {
	if x == nil {
		return
	}
	vertex, ok := x.(*Vertex)
	if !ok {
		panic(fmt.Sprintf("expected *model.Vertex, got %T", x))
	}
	*vs = append(*vs, vertex)
}

func (vs *Vertexes) Pop() interface{} {
	if len(*vs) == 0 {
		return nil
	}
	vertex := (*vs)[len(*vs)-1]
	*vs = (*vs)[:len(*vs)-1]
	return vertex
}

type Graph struct {
	VertexMapForSort  map[string]*Vertex
	VertexListForSort []*Vertex
	VertexMap         map[string]Vertex // for compatibility with old code
}

func IsDAG(g *Graph) error {
	if g == nil {
		return fmt.Errorf("graph is nil")
	}
	// must one vertex with in-degree 0
	if len(g.VertexMapForSort) == 0 {
		return fmt.Errorf("graph empty")
	}

	initialVertices := Vertexes{}
	for _, vertex := range g.VertexMapForSort {
		if vertex.InDegree == 0 {
			initialVertices = append(initialVertices, vertex)
		}
	}
	if len(initialVertices) == 0 {
		return fmt.Errorf("graph invalide, cycle exists in this graph")
	}
	heap.Init(&initialVertices)
	count := 0
	for len(initialVertices) > 0 {
		vertex := heap.Pop(&initialVertices).(*Vertex)
		if vertex.InDegree != 0 {
			return fmt.Errorf("error in topology sort, vertex %s has non-zero in-degree: %d", vertex.Name, vertex.InDegree)
		}
		count++
		for next := range vertex.Next {
			nextVertex, exists := g.VertexMapForSort[next]
			if !exists {
				continue
			}
			nextVertex.InDegree--
			if nextVertex.InDegree == 0 {
				heap.Push(&initialVertices, nextVertex)
			}
		}
	}
	if count != len(g.VertexMapForSort) {
		return fmt.Errorf("graph invalide, cycle exists in this graph")
	}
	return nil
}

func Extract(g *model.Graph) *Graph {
	if g == nil {
		return nil
	}

	graph := &Graph{
		VertexMapForSort: make(map[string]*Vertex),
		VertexMap:        make(map[string]Vertex),
	}
	for _, vertex := range g.Edges {
		if _, exists := graph.VertexMapForSort[vertex.From]; !exists {
			graph.VertexMapForSort[vertex.From] = &Vertex{
				Name:     vertex.From,
				InDegree: 0,
				Next:     make(map[string]struct{}),
			}
			graph.VertexListForSort = append(graph.VertexListForSort, graph.VertexMapForSort[vertex.From])
		}
		if _, exists := graph.VertexMapForSort[vertex.To]; !exists {
			graph.VertexMapForSort[vertex.To] = &Vertex{
				Name:     vertex.To,
				InDegree: 0,
				Next:     make(map[string]struct{}),
			}
			graph.VertexListForSort = append(graph.VertexListForSort, graph.VertexMapForSort[vertex.To])
		}
		graph.VertexMapForSort[vertex.To].InDegree++                     // Increment in-degree for the destination vertex
		graph.VertexMapForSort[vertex.From].Next[vertex.To] = struct{}{} // Add the destination vertex to the next map of the source vertex
	}

	for name, vertex := range graph.VertexMapForSort {
		graph.VertexMap[name] = *vertex
	}
	return graph
}
