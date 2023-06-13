package dispatcher

func (d dispatcher) registerInputs() {
	for _, input := range d.workflow.Inputs {
		input.Subscribe(d.dispatch)
	}
}
