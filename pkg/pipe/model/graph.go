package model

const (
	NebulaPipeExecBasicTag  = "pipe_exec_basic_tag"
	NebulaPipeExecBasicEdge = "pipe_exec_basic_edge"
)

type Vertex struct {
	Name        string `json:"name" yaml:"name" nebula:"name"`
	UUID        string `json:"uuid" yaml:"uuid" nebula:"uuid"`
	RuntimeUUID string `json:"runtime_uuid" yaml:"runtime_uuid" nebula:"runtime_uuid"`
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
	From string `json:"from" yaml:"from" nebula:"src_name"`
	To   string `json:"to" yaml:"to" nebula:"dst_name"`

	SrcUUID string `json:"src_uuid" yaml:"src_uuid" nebula:"src_uuid"`
	DstUUID string `json:"dst_uuid" yaml:"dst_uuid" nebula:"dst_uuid"`
}

func (v Edge) EdgeType() string {
	return NebulaPipeExecBasicEdge
}

func (v Edge) Dst() string {
	return v.DstUUID
}

func (v Edge) Src() string {
	return v.SrcUUID
}

func (v Edge) Props() map[string]interface{} {
	return map[string]interface{}{
		"src_name": v.From,
		"dst_name": v.To,
		"src_uuid": v.SrcUUID,
		"dst_uuid": v.DstUUID,
	}
}
