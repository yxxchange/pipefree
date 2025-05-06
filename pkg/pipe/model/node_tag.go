package model

import "encoding/json"

type NodeBasicTag struct {
	ApiVersion string `json:"apiVersion" norm:"api_version"`
	Kind       string `json:"kind" norm:"kind"`
	MetaData   string `json:"metadata" norm:"metadata"`
	Spec       string `json:"spec" norm:"spec"`
}

func (n NodeBasicTag) ConvertToNodeInfo() NodeInfo {
	var res NodeInfo
	res.ApiVersion = n.ApiVersion
	res.Kind = Kind(n.Kind)
	var spec Spec
	var meta MetaData
	_ = json.Unmarshal([]byte(n.Spec), &spec)
	_ = json.Unmarshal([]byte(n.MetaData), &meta)
	res.Spec = spec
	res.MetaData = meta
	return res
}
