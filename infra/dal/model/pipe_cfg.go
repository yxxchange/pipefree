package model

import (
	"encoding/json"
	"fmt"
)

const (
	EnvVarScopeGlobal = "global" // 全局环境变量
	EnvVarScopeNode   = "node"   // 节点环境变量
)

type PipeCfg struct {
	Basic
	Name      string  `json:"name" gorm:"column:name"`           // 流水线名称
	Namespace string  `json:"namespace" gorm:"column:namespace"` // 流水线命名空间
	Desc      string  `json:"desc" gorm:"column:desc"`           // 流水线描述
	EnvVars   EnvVars `json:"env_vars" gorm:"column:env_vars"`   // 流水线环境变量
	Graph     Graph   `json:"graph" gorm:"column:graph"`         // 流水线图结构
}

type EnvVars []EnvVar

func (envs *EnvVars) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), envs)
	case []byte:
		err = json.Unmarshal(v, envs)
	default:
		err = fmt.Errorf("unsupported type for EnvVars: %T", value)
	}
	return err
}

type EnvVar struct {
	Scope  string      `json:"scope"`  // 环境变量属性,如：全局、节点等
	Target string      `json:"target"` // 环境变量目标,如：节点名称
	Key    string      `json:"key"`    // 环境变量键
	Value  interface{} `json:"value"`  // 环境变量值
}

func (e *EnvVar) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), e)
	case []byte:
		err = json.Unmarshal(v, e)
	default:
		err = fmt.Errorf("unsupported type for EnvVar: %T", value)
	}
	return err
}

type Graph struct {
	Vertexes []Vertex `json:"vertexes"` // 图的顶点集合
	Edges    []Edge   `json:"edges"`    // 图的边集合
}

func (g *Graph) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), g)
	case []byte:
		err = json.Unmarshal(v, g)
	default:
		err = fmt.Errorf("unsupported type for Graph: %T", value)
	}
	return err
}

type Vertex struct {
	Name string `json:"name"` // 节点名称
}

type Edge struct {
	From string `json:"from"` // 边的起点
	To   string `json:"to"`   // 边的终点
}
