package model

import "encoding/json"

type NodeCfg struct {
	Basic
	Name string          `json:"name" gorm:"column:name"` // 节点名称
	Desc string          `json:"desc" gorm:"column:desc"` // 节点描述
	Spec json.RawMessage `json:"spec" gorm:"column:spec"` // 节点配置参数
}
