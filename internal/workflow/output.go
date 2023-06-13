package workflow

import (
	"encoding/json"
	"fmt"
	"github.com/valensto/ostraka/internal/workflow/provider/mqtt"
	"github.com/valensto/ostraka/internal/workflow/provider/sse"
)

type Publisher interface {
	Publish(events []byte)
}

type Output struct {
	Name        string
	Destination string
	Condition   *Condition
	Encoder     *Encoder

	Publisher Publisher
}

func UnmarshallOutput(name, destination string, condition *Condition, params any, opts Options) (*Output, error) {
	if name == "" {
		return nil, fmt.Errorf("output name is empty")
	}

	publisher, err := newPublisher(destination, params, opts)
	if err != nil {
		return nil, err
	}

	return &Output{
		Name:        name,
		Destination: destination,
		Condition:   condition,
		Encoder: &Encoder{
			format: JSON,
		},

		Publisher: publisher,
	}, nil
}

func newPublisher(dst string, params any, opts Options) (Publisher, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling input params: %w", err)
	}

	switch dst {
	case sse.SSE:
		return sse.NewPublisher(b, opts.Server, opts.Middlewares)

	case mqtt.MQTT:
		return mqtt.NewPublisher(b)

	default:
		return nil, fmt.Errorf("unknown publisher type: %s", dst)
	}
}

func (o *Output) Publish(event Event) error {
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
