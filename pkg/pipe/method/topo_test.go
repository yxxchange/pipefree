package method

import (
	"fmt"
	"github.com/yxxchange/pipefree/pkg/pipe/data"
	"testing"
)

func Test_TopologicalSort(t *testing.T) {
	graph1 := data.Graph{
		Edges: []data.Edge{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
			{From: "B", To: "D"},
			{From: "C", To: "D"},
			{From: "D", To: "E"},
			{From: "E", To: "F"},
		},
		Nodes: []data.Node{
			{MetaData: data.MetaData{Name: "A"}},
			{MetaData: data.MetaData{Name: "B"}},
			{MetaData: data.MetaData{Name: "C"}},
			{MetaData: data.MetaData{Name: "D"}},
			{MetaData: data.MetaData{Name: "E"}},
			{MetaData: data.MetaData{Name: "F"}},
			{MetaData: data.MetaData{Name: "G"}},
		},
	}
	graph2 := data.Graph{
		Edges: []data.Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
			{From: "C", To: "D"},
			{From: "D", To: "E"},
			{From: "E", To: "A"},
		},
		Nodes: []data.Node{
			{MetaData: data.MetaData{Name: "A"}},
			{MetaData: data.MetaData{Name: "B"}},
			{MetaData: data.MetaData{Name: "C"}},
			{MetaData: data.MetaData{Name: "D"}},
			{MetaData: data.MetaData{Name: "E"}},
		},
	}
	graph3 := data.Graph{
		Edges: []data.Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
			{From: "C", To: "D"},
			{From: "D", To: "E"},
			{From: "E", To: "F"},
			{From: "F", To: "C"},
		},
		Nodes: []data.Node{
			{MetaData: data.MetaData{Name: "A"}},
			{MetaData: data.MetaData{Name: "B"}},
			{MetaData: data.MetaData{Name: "C"}},
			{MetaData: data.MetaData{Name: "D"}},
			{MetaData: data.MetaData{Name: "E"}},
			{MetaData: data.MetaData{Name: "F"}},
			{MetaData: data.MetaData{Name: "G"}},
		},
	}
	graph4 := data.Graph{}

	graphList := []data.Graph{graph1, graph2, graph3, graph4}
	for _, graph := range graphList {
		sorter, err := NewTopologySorter().ExtractGraph(graph).TopologySort()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		for len(sorter.List) > 0 {
			nodes := sorter.GetZeroNode()
			for _, node := range nodes {
				fmt.Printf("node %s completed -- ", node.Node.Name)
			}
			fmt.Println()
		}
	}
}
