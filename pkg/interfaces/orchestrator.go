package interfaces

type Action string

const (
	ActionStart Action = "start"
	ActionPause Action = "pause"
	ActionStop  Action = "stop"
	ActionSkip  Action = "skip"
)

type Orchestrator[T any] interface {
	Serializer[T]
	// Launch the orchestrator
	// Launch() (T, error)
}
