package operator

import (
	"fmt"
	"github.com/yxxchange/pipefree/service/pipe_exec"
)

type WatchReq struct {
	Namespace    string `uri:"namespace" json:"namespace" form:"namespace" binding:"required"` // 命名空间
	Kind         string `uri:"kind" json:"kind" form:"kind" binding:"required"`                // 节点类型
	Version      string `uri:"version" json:"version" form:"version" binding:"required"`       // 节点版本
	OperatorName string `uri:"name" json:"name" form:"name"`                                   // 操作员名称
}

func (req *WatchReq) KeyPrefix() string {
	return fmt.Sprintf(pipe_exec.KeyPrefixTemplate, req.Namespace, req.Kind, req.Version)
}
