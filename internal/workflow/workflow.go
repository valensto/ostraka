package workflow

import "fmt"

type Workflow struct {
	Name    string
	Inputs  []*Input
	Outputs []*Output
}

func New(name string, inputs []*Input, outputs []*Output) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name is empty")
	}

	return &Workflow{
		Name:    name,
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}
