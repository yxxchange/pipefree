package model

const (
	BasicTagName  = "pipe_node_basic_tag"
	BasicEdgeName = "pipe_node_basic_edge"
)

type Vertex struct {
	Name        string `json:"name" yaml:"name"`
	UUID        string `json:"uuid" yaml:"uuid"` // static uuid
	RuntimeUUID string `json:"runtime_uuid" yaml:"runtime_uuid"`
}

func (v Vertex) VertexTagName() string {
	return BasicTagName
}

func (v Vertex) VertexID() string {
	return v.RuntimeUUID
}

type Edge struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`

	SrcUUID string `json:"src_uuid" yaml:"src_uuid" norm:"src_uuid"`
	DstUUID string `json:"dst_uuid" yaml:"dst_uuid" norm:"dst_uuid"`
}

func (v Edge) EdgeTypeName() string {
	return BasicEdgeName
}
