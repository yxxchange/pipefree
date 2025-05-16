package model

const (
	NebulaPipeExecBasicTag  = "pipe_exec_basic_tag"
	NebulaPipeExecBasicEdge = "pipe_exec_basic_edge"
)

type Vertex struct {
	Name        string `json:"name" yaml:"name" norm:"prop:name"`
	UUID        string `json:"uuid" yaml:"uuid" norm:"prop:uuid"`
	RuntimeUUID string `json:"runtime_uuid" yaml:"runtime_uuid" norm:"prop:runtime_uuid"`
}

func (v Vertex) TagName() string {
	return NebulaPipeExecBasicTag
}

func (v Vertex) VID() string {
	return v.RuntimeUUID
}

func (v Vertex) Props() map[string]interface{} {
	return map[string]interface{}{
		"runtime_uuid": v.RuntimeUUID,
		"name":         v.Name,
		"uuid":         v.UUID,
	}
}

type Edge struct {
	From string `json:"from" yaml:"from" norm:"prop:src_name"`
	To   string `json:"to" yaml:"to" norm:"prop:dst_name"`

	SrcUUID string `json:"src_uuid" yaml:"src_uuid" norm:"edge_src_id"`
	DstUUID string `json:"dst_uuid" yaml:"dst_uuid" norm:"edge_dst_id"`
}

func (v Edge) EdgeTypeName() string {
	return NebulaPipeExecBasicEdge
}
