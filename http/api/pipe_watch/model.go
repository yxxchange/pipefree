package pipe_watch

import (
	"fmt"
	"github.com/yxxchange/pipefree/service/pipe_exec"
)

type WatchReq struct {
	Namespace string `uri:"namespace" json:"namespace" form:"namespace" binding:"required"` // 命名空间
	Kind      string `uri:"kind" json:"kind" form:"kind" binding:"required"`                // 节点类型
}

func (req *WatchReq) KeyPrefix() string {
	return fmt.Sprintf(pipe_exec.KeyPrefixTemplate, req.Namespace, req.Kind)
}
