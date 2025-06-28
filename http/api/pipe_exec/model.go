package pipe_exec

import "github.com/yxxchange/pipefree/infra/dal/model"

type PipeExecView struct {
	Exec         model.PipeExec   `yaml:"pipe" json:"pipe"`   // 流水线执行信息
	NodeExecList []model.NodeExec `yaml:"nodes" json:"nodes"` // 节点执行列表
}

type PipeExecReqParam struct {
	PipeExecId int64 `uri:"pipe_exec_id" json:"pipe_exec_id" form:"pipe_exec_id"` // 流水线执行ID
	PipeId     int64 `uri:"pipe_id" json:"pipe_id" form:"pipe_id"`                // 流水线ID
}
