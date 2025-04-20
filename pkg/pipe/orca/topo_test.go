package orca

import (
	"fmt"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"testing"
)

func Test_TopologicalSort(t *testing.T) {
	graph1 := model.Graph{
		Edges: []model.Edge{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
			{From: "B", To: "D"},
			{From: "C", To: "D"},
			{From: "D", To: "E"},
			{From: "E", To: "F"},
		},
		Nodes: []model.Node{
			{MetaData: model.MetaData{Name: "A"}},
			{MetaData: model.MetaData{Name: "B"}},
			{MetaData: model.MetaData{Name: "C"}},
			{MetaData: model.MetaData{Name: "D"}},
			{MetaData: model.MetaData{Name: "E"}},
			{MetaData: model.MetaData{Name: "F"}},
			{MetaData: model.MetaData{Name: "G"}},
		},
	}
	graph2 := model.Graph{
		Edges: []model.Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
			{From: "C", To: "D"},
			{From: "D", To: "E"},
			{From: "E", To: "A"},
		},
		Nodes: []model.Node{
			{MetaData: model.MetaData{Name: "A"}},
			{MetaData: model.MetaData{Name: "B"}},
			{MetaData: model.MetaData{Name: "C"}},
			{MetaData: model.MetaData{Name: "D"}},
			{MetaData: model.MetaData{Name: "E"}},
		},
	}
	graph3 := model.Graph{
		Edges: []model.Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
			{From: "C", To: "D"},
			{From: "D", To: "E"},
			{From: "E", To: "F"},
			{From: "F", To: "C"},
		},
		Nodes: []model.Node{
			{MetaData: model.MetaData{Name: "A"}},
			{MetaData: model.MetaData{Name: "B"}},
			{MetaData: model.MetaData{Name: "C"}},
			{MetaData: model.MetaData{Name: "D"}},
			{MetaData: model.MetaData{Name: "E"}},
			{MetaData: model.MetaData{Name: "F"}},
			{MetaData: model.MetaData{Name: "G"}},
		},
	}
	graph4 := model.Graph{}

	graphList := []model.Graph{graph1, graph2, graph3, graph4}
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
