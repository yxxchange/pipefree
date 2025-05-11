package model

type PipeConfig struct {
	PipeFlow `json:",inline" yaml:",inline"`
}

func (p PipeConfig) ToPipeExec() PipeExec {
	return p.ToExec()
}

type PipeExec struct {
	PipeFlow `json:",inline" yaml:",inline"`
}
