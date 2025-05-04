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
