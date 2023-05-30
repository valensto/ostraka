package dispatcher

import (
	"fmt"

	"github.com/valensto/ostraka/internal/input/mqtt"
	"github.com/valensto/ostraka/internal/input/webhook"
	"github.com/valensto/ostraka/internal/workflow"
)

type InputProvider interface {
	Subscribe() error
}

func (d dispatcher) subscribeInputs() error {
	for _, i := range d.workflow.Inputs {
		provider, err := d.getInputProvider(i)
		if err != nil {
			return fmt.Errorf("error getting input provider: %w", err)
		}

		err = provider.Subscribe()
		if err != nil {
			return fmt.Errorf("error subscribing input: %w", err)
		}
	}

	return nil
}

func (d dispatcher) getInputProvider(i workflow.Input) (InputProvider, error) {
	switch i.Type {
	case workflow.Webhook:
		return webhook.New(i, d.router, d.inputEvents)
	case workflow.MQTT:
		return mqtt.New(i, d.inputEvents)
	default:
		return nil, fmt.Errorf("unknown input type: %s", i.Type)
	}
}
