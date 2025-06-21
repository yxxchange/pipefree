package model

import "encoding/json"

type NodeCfg struct {
	Basic
	Name      string          `json:"name" gorm:"column:name"`               // 节点名称
	Desc      string          `json:"desc" gorm:"column:desc"`               // 节点描述
	PipeCfgId int64           `json:"pipe_cfg_id" gorm:"column:pipe_cfg_id"` // 流水线配置ID
	InDegree  int             `json:"in_degree" gorm:"column:in_degree"`     // 入度
	Spec      json.RawMessage `json:"spec" gorm:"column:spec"`               // 节点配置参数
}

func (*NodeCfg) TableName() string {
	return "node_cfg"
}
