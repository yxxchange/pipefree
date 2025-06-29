package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	EnvVarScopeGlobal = "global" // 全局环境变量
	EnvVarScopeNode   = "node"   // 节点环境变量
)

type PipeCfg struct {
	Basic
	Name    string   `json:"name" yaml:"name" gorm:"column:name"`            // 流水线名称
	Space   string   `json:"space" yaml:"space" gorm:"column:space"`         // 项目空间
	Desc    string   `json:"desc" yaml:"desc" gorm:"column:desc"`            // 流水线描述
	Version int      `json:"version" yaml:"version" gorm:"column:version"`   // 流水线版本
	EnvVars *EnvVars `json:"env_vars" yaml:"envVars" gorm:"column:env_vars"` // 流水线环境变量
	Graph   *Graph   `json:"graph" yaml:"graph" gorm:"column:graph"`         // 流水线图结构
}

func (*PipeCfg) TableName() string {
	return "pipe_cfg"
}

func (p *PipeCfg) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), p)
	case []byte:
		err = json.Unmarshal(v, p)
	default:
		err = fmt.Errorf("unsupported type for Config: %T", value)
	}
	return err
}

func (p *PipeCfg) Value() (value driver.Value, err error) {
	if p == nil {
		return nil, nil
	}
	value, err = json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PipeCfg: %w", err)
	}
	return value, nil
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

func (envs *EnvVars) Value() (value driver.Value, err error) {
	if envs == nil {
		return nil, nil
	}
	value, err = json.Marshal(envs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal EnvVars: %w", err)
	}
	return value, nil
}

type EnvVar struct {
	Scope  string      `json:"scope" yaml:"scope"`   // 环境变量属性,如：全局、节点等
	Target string      `json:"target" yaml:"target"` // 环境变量目标,如：节点名称
	Key    string      `json:"key" yaml:"key"`       // 环境变量键
	Value  interface{} `json:"value" yaml:"value"`   // 环境变量值
}

type Graph struct {
	Edges []Edge `json:"edges" yaml:"edges"` // 图的边集合
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

func (g *Graph) Value() (value driver.Value, err error) {
	if g == nil {
		return nil, nil
	}
	value, err = json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Graph: %w", err)
	}
	return value, nil
}

// Edge represents a directed edge in the graph.
type Edge struct {
	From string `json:"from" yaml:"from"` // 边的起点
	To   string `json:"to" yaml:"to"`     // 边的终点
}
