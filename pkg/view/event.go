package view

import "github.com/yxxchange/pipefree/pkg/pipe/data"

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

type Event struct {
	EventType EventType
	Kind      string
	data.Node
}

type NodeView struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	data.MetaData
	data.Spec
	data.Status
}
