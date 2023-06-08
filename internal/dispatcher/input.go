package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/webhook"

	"github.com/valensto/ostraka/internal/workflow"
)

type subscriber interface {
	Subscribe(dispatch func(input *workflow.Input, data []byte) error) error
}

func (d dispatcher) registerInputs() error {
	for _, input := range d.workflow.Inputs {
		s, err := d.getSubscriber(input)
		if err != nil {
			return fmt.Errorf("error getting subscriber: %w", err)
		}

		err = s.Subscribe(d.dispatch)
		if err != nil {
			return fmt.Errorf("error registering subscriber: %w", err)
		}
	}

	return nil
}

func (d dispatcher) getSubscriber(input *workflow.Input) (subscriber, error) {
	switch input.Source {
	case workflow.Webhook:
		return webhook.NewSubscriber(input, d.server)
	case workflow.MQTTSub:
		return mqtt.NewSubscriber(input)
	default:
		return nil, fmt.Errorf("unknown input type: %s", input.Source)
	}
}
