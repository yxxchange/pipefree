package pipe_perm

type PipePermReqParam struct {
	Space     string `uri:"space" json:"space" binding:"required"`         // 流水线空间
	Namespace string `uri:"namespace" json:"namespace" binding:"required"` // 节点命名空间
}
