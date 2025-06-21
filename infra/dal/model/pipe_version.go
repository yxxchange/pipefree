package model

type PipeVersion struct {
	Basic
	PipeCfgId int64   `json:"pipe_cfg_id" gorm:"column:pipe_cfg_id"`     // 流水线配置ID
	Version   int     `json:"version" gorm:"column:version"`             // 流水线版本号
	Config    PipeCfg `json:"pipe_cfg" gorm:"column:pipe_cfg;type:json"` // 流水线配置内容
}
