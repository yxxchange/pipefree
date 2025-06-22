package pipe_cfg

import "github.com/yxxchange/pipefree/infra/dal/model"

// PipeView 流水线视图模型
type PipeView struct {
	model.PipeCfg `yaml:",inline" json:",inline"` // 流水线配置
}

func Convert(cfg model.PipeCfg) PipeView {
	return PipeView{
		PipeCfg: cfg,
	}
}

type PipeReqParam struct {
	PipeId int64  `uri:"pipe_id" json:"pipe_id" form:"pipe_id"`
	Yaml   string `form:"yaml"`
}
