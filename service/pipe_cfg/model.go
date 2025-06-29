package pipe_cfg

import (
	"github.com/yxxchange/pipefree/infra/dal/model"
	"github.com/yxxchange/pipefree/service/graph"
	"github.com/yxxchange/pipefree/utils"
)

type PipeComponent struct {
	Pipe     *model.PipeCfg
	NodeList []*model.NodeCfg
	NodeMap  map[string]*model.NodeCfg

	Namespaces []string
	Graph      *graph.Graph
}

func NewPipeComponent(pipe *model.PipeCfg, nodes []*model.NodeCfg) *PipeComponent {
	nodeMap := make(map[string]*model.NodeCfg)
	nsList := make([]string, 0, len(nodes))
	for _, node := range nodes {
		nodeMap[node.Name] = node
		nsList = append(nsList, node.Namespace)
	}
	nsList = utils.Deduplicate(nsList)
	return &PipeComponent{
		Pipe:       pipe,
		NodeList:   nodes,
		NodeMap:    nodeMap,
		Namespaces: nsList,
		Graph:      graph.Extract(pipe.Graph),
	}
}
