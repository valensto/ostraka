package workflow

import (
	"fmt"
	"github.com/valensto/ostraka/internal/event"
	"github.com/valensto/ostraka/internal/provider"
)

type Output struct {
	Name      string
	Condition *Condition
	Encoder   *event.Encoder

	Publisher provider.Publisher
}

func UnmarshallOutput(name, dst string, condition *Condition, encoder *event.Encoder, params any, opts provider.Options) (*Output, error) {
	if name == "" {
		return nil, fmt.Errorf("output name is empty")
	}

	publisher, err := provider.NewPublisher(dst, params, opts)
	if err != nil {
		return nil, fmt.Errorf("error creating publisher for output %s got: %w", name, err)
	}

	return &Output{
		Name:      name,
		Condition: condition,
		Encoder:   encoder,

		Publisher: publisher,
	}, nil
}

func (o *Output) Publish(event event.Payload) error {
	if !o.Condition.Match(event) {
		return fmt.Errorf("event does not match output %s condition", o.Name)
	}

	b, err := o.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("error encoding event for output %s got: %w", o.Name, err)
	}

	o.Publisher.Publish(b)
	return nil
}
