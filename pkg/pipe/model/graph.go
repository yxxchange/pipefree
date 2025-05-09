package model

import "github.com/yxxchange/pipefree/helper/serialize"

const (
	BasicTagName  = "pipe_node_basic_tag"
	BasicEdgeName = "pipe_node_basic_edge"
)

type PipeFragment struct {
	Vertexes      []NodeBasicTag
	Edges         []BasicEdge
	NodeSnapshots []interface{}
}

type NodeBasicTag struct {
	VID        string `json:"vid" norm:"vid"`
	ApiVersion string `json:"apiVersion" norm:"api_version"`
	Kind       string `json:"kind" norm:"kind"`
	MetaData   string `json:"metadata" norm:"metadata"`
	Spec       string `json:"spec" norm:"spec"`
}

func (n NodeBasicTag) ToIdentifier() NodeIdentifier {
	var meta MetaData
	_ = serialize.JsonDeserialize([]byte(n.MetaData), &meta)
	return NodeIdentifier{
		ApiVersion: n.ApiVersion,
		Kind:       Kind(n.Kind),
		Namespace:  meta.Namespace,
		Operation:  meta.Operation,
	}
}

func (n NodeBasicTag) ToString() string {
	b, _ := serialize.JsonSerialize(n)
	return string(b)
}

func (n NodeBasicTag) VertexTagName() string {
	return BasicTagName
}

func (n NodeBasicTag) VertexID() string {
	return n.VID
}

type BasicEdge struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
}

func (e BasicEdge) EdgeTypeName() string {
	return BasicEdgeName
}

func (n NodeBasicTag) ConvertToNodeInfo() Node {
	var res Node
	res.ApiVersion = n.ApiVersion
	res.Kind = Kind(n.Kind)
	var spec Spec
	var meta MetaData
	_ = serialize.JsonDeserialize([]byte(n.Spec), &spec)
	_ = serialize.JsonDeserialize([]byte(n.MetaData), &meta)
	res.Spec = spec
	res.MetaData = meta
	return res
}
