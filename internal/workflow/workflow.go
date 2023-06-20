package workflow

import (
	"fmt"
	"github.com/valensto/ostraka/internal/middleware"
	"github.com/valensto/ostraka/internal/provider"
)

type Workflow struct {
	Name string `json:"name" yaml:"name" validate:"required"`
	Slug string `json:"-" yaml:"-"`

	EventType eventType `json:"event_type" yaml:"event_type" validate:"required,dive,required"`

	Middlewares *middleware.Middlewares `json:"middlewares" yaml:"middlewares"`

	Inputs  []*input  `json:"inputs" yaml:"inputs" validate:"required,dive,required"`
	Outputs []*output `json:"outputs" yaml:"outputs" validate:"required,dive,required"`

	consumers []consumer
}

func (wf *Workflow) Init(opts provider.Options) error {
	err := wf.loadInputs(opts)
	if err != nil {
		return fmt.Errorf("error unmarshalling inputs: %w", err)
	}

	err = wf.loadOutputs(opts)
	if err != nil {
		return fmt.Errorf("error unmarshalling outputs: %w", err)
	}

	return nil
}

func (wf *Workflow) dispatch(from *input, ib []byte) {
	collect := wf.collect(from, ib)
	defer collect.consumes()

	p, err := from.Decoder.decode(ib)
	if err != nil {
		collect.withError(err)
		return
	}

	for _, o := range wf.Outputs {
		ob, err := o.publish(p)
		if err != nil {
			collect.withError(err)
			continue
		}

		collect.addOutput(o, ob, err)
	}
}
