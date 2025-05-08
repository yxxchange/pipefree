package model

type PipeConfig struct {
	Node
}

type PipeExec struct {
	NodeSnapshot
}

type NodeSnapshot struct {
	VID string `json:"vid"`
	Node
	Graph *GraphSnapshot `json:"graph,omitempty" yaml:"graph,omitempty"`
}

type GraphSnapshot struct {
	Edges     []Edge         `json:"edges,omitempty" yaml:"edges,omitempty"`
	Vertexes  []NodeSnapshot `json:"vertexes,omitempty" yaml:"vertexes,omitempty"`
	Reference MetaData       `json:"reference,omitempty" yaml:"reference,omitempty"`
}

func (p PipeConfig) ToPipeExec() PipeExec {
	pipeExec := PipeExec{
		NodeSnapshot: p.Node.ToSnapshot(),
	}
	return pipeExec
}
