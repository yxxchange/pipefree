package operator

type WatchReq struct {
	Namespace string `uri:"namespace" json:"namespace" form:"namespace"` // 命名空间
	Kind      string `uri:"kind" json:"kind" form:"kind"`                // 节点类型
	Version   string `uri:"version" json:"version" form:"version"`       // 节点版本
}
