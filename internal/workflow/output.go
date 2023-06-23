package workflow

import (
	"fmt"
	"github.com/valensto/ostraka/internal/provider"
)

type output struct {
	Name        string     `json:"name" yaml:"name" validate:"required"`
	Destination string     `json:"destination" yaml:"destination" validate:"required"`
	Condition   *condition `json:"condition,omitempty" yaml:"condition,omitempty"`
	Encoder     *encoder   `json:"encoder" yaml:"encoder"`

	Params    any `json:"params" yaml:"params" validate:"required"`
	publisher provider.Publisher
}

func (o *output) loadPublisher(opts provider.Options) error {
	publisher, err := provider.NewPublisher(o.Destination, o.Params, opts)
	if err != nil {
		return err
	}

	o.publisher = publisher
	o.Params = nil
	return nil
}

func (wf *Workflow) loadOutputs(opts provider.Options) error {
	for i, _ := range wf.Outputs {
		if err := wf.Outputs[i].loadPublisher(opts); err != nil {
			return fmt.Errorf("error unmarshalling output %s got: %w", wf.Outputs[i].Name, err)
		}

		var c *condition
		if wf.Outputs[i].Condition != nil {
			uc, err := wf.Outputs[i].Condition.computeConditions()
			if err != nil {
				return fmt.Errorf("error converting condition yaml: %w", err)
			}

			c = uc
		}

		fmt.Printf("condition for output %s: %v\n", wf.Outputs[i].Name, c)
		wf.Outputs[i].Condition = c
	}

	return nil
}

func (o *output) publish(event payload) ([]byte, error) {
	if o.publisher == nil {
		return nil, fmt.Errorf("output %s is not initialized", o.Name)
	}

	if !o.Condition.match(event) {
		return nil, fmt.Errorf("event does not match output %s condition", o.Name)
	}

	b, err := o.Encoder.encode(event)
	if err != nil {
		return nil, fmt.Errorf("error encoding event for output %s got: %w", o.Name, err)
	}

	o.publisher.Publish(b)
	return b, nil
}
