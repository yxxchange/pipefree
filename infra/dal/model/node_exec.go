package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	NodePhaseReady     = "Ready"     // 节点就绪状态
	NodePhaseRunning   = "Running"   // 节点运行中状态
	NodePhaseSucceeded = "Succeeded" // 节点成功状态
	NodePhaseFailed    = "Failed"    // 节点失败状态
	NodePhasePending   = "Pending"   // 节点等待状态
	NodePhaseSkipped   = "Skipped"   // 节点跳过状态
	NodePhaseUnknown   = "Unknown"   // 节点未知状态
)

func NewNodeExec(node *NodeCfg, pipe *PipeExec) *NodeExec {
	return &NodeExec{
		Name:        node.Name,
		Namespace:   pipe.Namespace,
		PipeName:    pipe.Name,
		PipeVersion: pipe.Version,
		NodeCfgId:   node.Id,
		PipeCfgId:   node.PipeCfgId,
		PipeExecId:  pipe.Id,
		InDegree:    node.InDegree,
		Spec:        node.Spec,
		Phase: &NodePhase{
			Phase: NodePhaseReady,
		},
	}
}

type NodeExec struct {
	Basic
	Name        string     `json:"name" gorm:"column:name"`                 // 节点名称
	Namespace   string     `json:"namespace" gorm:"column:namespace"`       // 节点命名空间
	PipeName    string     `json:"pipe_name" gorm:"column:pipe_name"`       // 流水线名称
	PipeVersion int        `json:"pipe_version" gorm:"column:pipe_version"` // 流水线版本
	NodeCfgId   int64      `json:"node_cfg_id" gorm:"column:node_cfg_id"`   // 节点配置ID
	PipeCfgId   int64      `json:"pipe_cfg_id" gorm:"column:pipe_cfg_id"`   // 流水线配置ID
	PipeExecId  int64      `json:"pipe_exec_id" gorm:"column:pipe_exec_id"` // 流水线执行ID
	InDegree    int        `json:"inDegree" gorm:"column:in_degree"`        // 入度
	Spec        *Kv        `json:"spec" gorm:"column:spec"`                 // 节点执行参数
	Phase       *NodePhase `json:"status" gorm:"column:status"`             // 节点执行节点
}

func (*NodeExec) TableName() string {
	return "node_exec"
}

type NodePhase struct {
	Phase  string       `json:"phase"`  // 节点状态，如：Pending、Running、Succeeded、Failed等
	Chains []PhaseChain `json:"chains"` // 状态快照链
}

func (n *NodePhase) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), n)
	case []byte:
		err = json.Unmarshal(v, n)
	default:
		err = fmt.Errorf("unsupported type for NodePhase: %T", value)
	}
	return err
}

func (n *NodePhase) Value() (value driver.Value, err error) {
	if n == nil {
		return nil, nil
	}
	value, err = json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NodePhase: %w", err)
	}
	return value, nil
}

type PhaseChain struct {
	Time      time.Time       `json:"time"`      // 状态快照时间
	Phase     string          `json:"phase"`     // 节点状态，如：Pending、Running、Succeeded、Failed等
	Comment   string          `json:"comment"`   // 状态快照备注
	Customize json.RawMessage `json:"customize"` // 自定义状态快照内容
}
