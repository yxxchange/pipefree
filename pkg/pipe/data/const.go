package data

type Phase string
type NodeType string

const (
	PhaseReady          Phase = "ready"      // the ready phase of the node means the pipe is ready to run
	PhaseRunning        Phase = "running"    // the running phase of the node means the pipe is running
	PhasePipeSucceed    Phase = "succeed"    // the succeed phase of the node means the pipe is succeeded
	PhasePipeFailed     Phase = "failed"     // the failed phase of the node means the pipe is failed because of some inner error
	PhasePipeTerminated Phase = "terminated" // the terminated phase of the node means the pipe is terminated by the user
	PhasePipePaused     Phase = "paused"     // the paused phase of the node means the pipe is paused by the user

	NodeTypeScalar   NodeType = "scalar"   // the scalar node type means the node not contains any sub nodes
	NodeTypeCompound NodeType = "compound" // the compound node type means the node contains some sub nodes
)
