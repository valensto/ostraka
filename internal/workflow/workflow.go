package workflow

type Workflow struct {
	Inputs  map[string]Input
	Outputs map[string]Output
}

func New(inputs []*Input, outputs []*Output) (*Workflow, error) {
	wf := Workflow{
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
