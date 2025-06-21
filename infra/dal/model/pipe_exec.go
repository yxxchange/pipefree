package model

const (
	PipeExecStateUnknown int = -1 // 未知状态
	PipeExecStateRunning int = 0  // 流水线执行中
	PipeExecStateSuccess int = 1  // 流水线执行成功
	PipeExecStateFailed  int = 2  // 流水线执行失败
)

type PipeExec struct {
	Basic
	PipeCfgId int64   `json:"pipe_cfg_id" gorm:"column:pipe_cfg_id"`     // 流水线配置ID
	Version   int     `json:"version" gorm:"column:version"`             // 流水线版本号
	Graph     Graph   `json:"graph" gorm:"column:graph;type:json"`       // 流水线图内容
	EnvVars   EnvVars `json:"env_vars" gorm:"column:env_vars;type:json"` // 环境变量
	State     int     `json:"state" gorm:"column:state"`                 // 流水线执行状态
}
