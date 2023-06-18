package workflow

import (
	"fmt"
	"github.com/gosimple/slug"
)

type Workflow struct {
	Name string
	Slug string

	Inputs  []*Input
	Outputs []*Output

	consumers []consumer
}

func New(name string, inputs []*Input, output []*Output) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name is empty")
	}

	wf := Workflow{
		Name: name,
		Slug: slug.Make(name),

		Inputs:  inputs,
		Outputs: output,
	}

	return &wf, nil
}

func (wf *Workflow) Listen(consumers ...consumer) error {
	wf.consumers = consumers

	for _, input := range wf.Inputs {
		err := input.listen(wf.dispatch)
		if err != nil {
			return fmt.Errorf("error subscribing input %s got: %w", input.Name, err)
		}
	}

	return nil
}

func (wf *Workflow) dispatch(from *Input, bytes []byte) {
	collect := wf.collect(from, bytes)
	defer collect.consumes()

	payload, err := from.Decoder.Decode(bytes)
	if err != nil {
		collect.withError(err)
		return
	}

	for _, output := range wf.Outputs {
		b, err := output.Publish(payload)
		collect.addOutput(output, b, err)
	}
}
