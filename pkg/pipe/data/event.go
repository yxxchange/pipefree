package data

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

type Event struct {
	EventType EventType
	Kind      string
	Node
}

type NodeView struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	MetaData
	Spec
	Status
}
