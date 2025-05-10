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

func (n NodeSnapshot) toGraphMeta() GraphMeta {
	var res GraphMeta
	res.Vertexes = []NodeBasicTag{n.ToBasicTag()}
	if n.Graph != nil {
		for _, edge := range n.Graph.Edges {
			res.Edges = append(res.Edges, BasicEdge{
				SrcID: edge.From,
				DstID: edge.To,
			})
		}
	}
	for _, vertex := range n.Graph.Vertexes {
		son := vertex.toGraphMeta()
		res.Vertexes = append(res.Vertexes, son.Vertexes...)
		res.Edges = append(res.Edges, son.Edges...)
	}
	return res
}

type GraphSnapshot struct {
	Edges     []Edge         `json:"edges,omitempty" yaml:"edges,omitempty"`
	Vertexes  []NodeSnapshot `json:"vertexes,omitempty" yaml:"vertexes,omitempty"`
	Reference MetaData       `json:"reference,omitempty" yaml:"reference,omitempty"`
}

func (p PipeExec) ToGraph() GraphMeta {
	return p.NodeSnapshot.toGraphMeta()
}
