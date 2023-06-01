package workflow

import "fmt"

type Workflow struct {
	Name    string
	Inputs  map[string]Input
	Outputs map[string]Output
}

func New(name string, inputs []*Input, outputs []*Output) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name is empty")
	}

	wf := Workflow{
		Name:    name,
		Inputs:  make(map[string]Input),
		Outputs: make(map[string]Output),
	}

	for _, input := range inputs {
		wf.Inputs[input.Name] = *input
	}

	for _, output := range outputs {
		wf.Outputs[output.Name] = *output
	}

	return &wf, nil
}
