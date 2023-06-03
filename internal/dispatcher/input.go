package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/webhook"

	"github.com/valensto/ostraka/internal/workflow"
)

type Subscriber interface {
	Subscribe(events chan<- map[string]any) error
}

func (d dispatcher) subscribeInputs() error {
	for _, input := range d.workflow.Inputs {
		subscriber, err := d.getInputProvider(input)
		if err != nil {
			return fmt.Errorf("error getting input subscriber: %w", err)
		}

		err = subscriber.Subscribe(d.inputEvents)
		if err != nil {
			return fmt.Errorf("error subscribing input: %w", err)
		}
	}

	return nil
}

func (d dispatcher) getInputProvider(input workflow.Input) (Subscriber, error) {
	switch input.Source {
	case workflow.Webhook:
		return webhook.New(input, d.server)
	case workflow.MQTTSub:
		return mqtt.NewSubscriber(input)
	default:
		return nil, fmt.Errorf("unknown input type: %s", input.Source)
	}
}
