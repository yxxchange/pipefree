package graph

import (
	"container/heap"
	"fmt"
	"github.com/yxxchange/pipefree/infra/dal/model"
)

type Vertexes []*model.Vertex

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
	vertex, ok := x.(*model.Vertex)
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

type DAG struct {
	VertexMapForSort map[string]*model.Vertex
	VertexMap        map[string]model.Vertex // for compatibility with old code
}

func IsDAG(dag *DAG) error {
	// must one vertex with in-degree 0
	if len(dag.VertexMapForSort) == 0 {
		return fmt.Errorf("graph empty")
	}

	initialVertices := Vertexes{}
	for _, vertex := range dag.VertexMapForSort {
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
		vertex := heap.Pop(&initialVertices).(*model.Vertex)
		if vertex.InDegree != 0 {
			return fmt.Errorf("error in topology sort, vertex %s has non-zero in-degree: %d", vertex.Name, vertex.InDegree)
		}
		count++
		for _, next := range vertex.Next {
			nextVertex, exists := dag.VertexMapForSort[next]
			if !exists {
				continue
			}
			nextVertex.InDegree--
			if nextVertex.InDegree == 0 {
				heap.Push(&initialVertices, nextVertex)
			}
		}
	}
	if count != len(dag.VertexMapForSort) {
		return fmt.Errorf("graph invalide, cycle exists in this graph")
	}
	return nil
}

func Extract(g *model.Graph) (*DAG, error) {
	dag := &DAG{
		VertexMapForSort: make(map[string]*model.Vertex),
		VertexMap:        make(map[string]model.Vertex),
	}
	if g == nil {
		return nil, fmt.Errorf("graph is empty")
	}
	for i, vertex := range g.Vertexes {
		if _, exists := dag.VertexMapForSort[vertex.Name]; exists {
			return nil, fmt.Errorf("graph invalide, duplicate vertex name: %s", vertex.Name)
		}
		dag.VertexMapForSort[vertex.Name] = &g.Vertexes[i]
	}
	for _, edge := range g.Edges {
		if _, exists := dag.VertexMapForSort[edge.From]; !exists {
			return nil, fmt.Errorf("graph invalide, edge from vertex not found: %s", edge.From)
		}
		if _, exists := dag.VertexMapForSort[edge.To]; !exists {
			return nil, fmt.Errorf("graph invalide, edge to vertex not found: %s", edge.To)
		}
		if edge.From == edge.To {
			return nil, fmt.Errorf("graph invalide, self-loop detected at vertex: %s", edge.From)
		}
		dag.VertexMapForSort[edge.To].InDegree++
		dag.VertexMapForSort[edge.From].Next = append(dag.VertexMapForSort[edge.From].Next, edge.To)
	}

	for name, vertex := range dag.VertexMapForSort {
		dag.VertexMap[name] = *vertex
	}
	return dag, nil
}
