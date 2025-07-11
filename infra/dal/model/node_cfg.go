package model

import (
	"database/sql/driver"
	"encoding/json"
)

type NodeCfg struct {
	Basic
	Name      string `json:"name" yaml:"name" gorm:"column:name"`                      // 节点名称
	Kind      string `json:"kind" yaml:"kind" gorm:"column:kind"`                      // 节点类型
	Version   string `json:"version" yaml:"version" gorm:"column:version"`             // 节点版本
	Namespace string `json:"namespace" yaml:"namespace" gorm:"column:namespace"`       // 节点命名空间
	PipeSpace string `json:"pipe_space" yaml:"pipe_space" gorm:"column:pipe_space"`    // 流水线命名空间
	PipeName  string `json:"pipe_name" yaml:"pipe_name" gorm:"column:pipe_name"`       // 流水线名称
	Desc      string `json:"desc" yaml:"desc" gorm:"column:desc"`                      // 节点描述
	PipeCfgId int64  `json:"pipe_cfg_id" yaml:"pipe_cfg_id" gorm:"column:pipe_cfg_id"` // 流水线配置ID
	InDegree  int    `json:"in_degree" yaml:"in_degree" gorm:"column:in_degree"`       // 入度
	Spec      *Kv    `json:"spec" yaml:"spec" gorm:"column:spec"`                      // 节点配置参数
}

func (*NodeCfg) TableName() string {
	return "node_cfg"
}

type Kv map[string]interface{}

func (kv *Kv) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), kv)
	case []byte:
		err = json.Unmarshal(v, kv)
	default:
		err = nil
	}
	return err
}

func (kv *Kv) Value() (value driver.Value, err error) {
	if kv == nil {
		return nil, nil
	}
	value, err = json.Marshal(kv)
	if err != nil {
		return nil, err
	}
	return value, nil
}
