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

func New(name string, inputs []*Input, output []*Output, consumers ...consumer) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name is empty")
	}

	wf := Workflow{
		Name: name,
		Slug: slug.Make(name),

		Inputs:  inputs,
		Outputs: output,

		consumers: consumers,
	}

	return &wf, nil
}

func (wf *Workflow) Listen() error {
	for _, input := range wf.Inputs {
		err := input.listen(wf.dispatch)
		if err != nil {
			return fmt.Errorf("error subscribing input %s got: %w", input.Name, err)
		}
	}

	return nil
}

func (wf *Workflow) dispatch(from *Input, data []byte) {
	c := wf.newCollector(from, data)
	defer c.dump()

	e, err := from.Decoder.Decode(data)
	if err != nil {
		c.dump()
		return
	}

	for _, output := range wf.Outputs {
		err := output.Publish(e)
		if err != nil {
			c.addError(err)
		}
	}
}
