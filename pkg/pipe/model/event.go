package model

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

type Event struct {
	EventType EventType
	Data      []byte
}

type NodeView struct {
	ApiVersion string `json:"apiVersion"`
	Kind       `json:"kind" yaml:"kind"`
	MetaData   `json:"metaData" yaml:"metaData"`
	Spec       `json:"spec" yaml:"spec"`
	Status     `json:"status" yaml:"status"`
}
