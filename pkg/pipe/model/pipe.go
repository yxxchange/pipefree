package model

type PipeConfig struct {
	PipeFlow
}

func (p PipeConfig) ToPipeExec() PipeExec {
	return p.ToExec()
}

type PipeExec struct {
	PipeFlow
}
