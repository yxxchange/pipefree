package pipe_cfg

import "github.com/yxxchange/pipefree/infra/dal/model"

// PipeView 流水线视图模型
type PipeView struct {
	PipeCfg     model.PipeCfg   `yaml:"pipe" json:"pipe"` // 流水线配置
	NodeCfgList []model.NodeCfg `yaml:"nodes" json:"nodes"`
}

func Convert(cfg model.PipeCfg) PipeView {
	return PipeView{
		PipeCfg: cfg,
	}
}

type PipeReqParam struct {
	PipeId int64    `uri:"pipe_id" json:"pipe_id" form:"pipe_id"`
	View   PipeView `yaml:"view" json:"view"`
}
