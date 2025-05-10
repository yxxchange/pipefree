package model

type PipeConfig struct {
	Node
}

func (p PipeConfig) ToPipeExec() PipeExec {
	pipeExec := PipeExec{
		NodeSnapshot: p.Node.ToSnapshot(),
	}
	return pipeExec
}

type PipeExec struct {
	NodeSnapshot
}

type NodeSnapshot struct {
	VID string `json:"vid"`
	Node
	Graph *GraphSnapshot `json:"graph,omitempty" yaml:"graph,omitempty"`
}

func (n NodeSnapshot) decompose() PipeFragment {
	var res PipeFragment
	res.Vertexes = []NodeBasicTag{n.ToBasicTag()}
	res.NodeSnapshots = []interface{}{n}
	if n.Graph != nil {
		for _, edge := range n.Graph.Edges {
			res.Edges = append(res.Edges, BasicEdge{
				SrcID: edge.From,
				DstID: edge.To,
			})
		}
	}
	for _, vertex := range n.Graph.Vertexes {
		son := vertex.decompose()
		res.Vertexes = append(res.Vertexes, son.Vertexes...)
		res.Edges = append(res.Edges, son.Edges...)
		res.NodeSnapshots = append(res.NodeSnapshots, son.NodeSnapshots...)
	}
	return res
}

type GraphSnapshot struct {
	Edges     []Edge         `json:"edges,omitempty" yaml:"edges,omitempty"`
	Vertexes  []NodeSnapshot `json:"vertexes,omitempty" yaml:"vertexes,omitempty"`
	Reference Reference      `json:"reference,omitempty" yaml:"reference,omitempty"`
}

func (p PipeExec) Decompose() PipeFragment {
	return p.NodeSnapshot.decompose()
}
