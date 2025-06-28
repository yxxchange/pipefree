package operator

type ListAndWatchReq struct {
	Namespace string `uri:"namespace" json:"namespace" form:"namespace"` // 命名空间
	Name      string `uri:"name" json:"name" form:"name"`                // 流水线名称
}
