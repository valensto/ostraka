package dispatcher

import (
	"fmt"

	"github.com/valensto/ostraka/internal/input/mqtt"
	"github.com/valensto/ostraka/internal/input/webhook"
	"github.com/valensto/ostraka/internal/workflow"
)

func (d dispatcher) subscribeInputs() error {
	for _, i := range d.workflow.Inputs {
		switch i.Type {
		case workflow.Webhook:
			input, err := webhook.New(i, d.router, d.inputEvents)
			if err != nil {
				return fmt.Errorf("error creating webhook input: %w", err)
			}
			err = input.Subscribe()
			if err != nil {
				return fmt.Errorf("error subscribing webhook input: %w", err)
			}

		case workflow.MQTT:
			input, err := mqtt.New(i, d.inputEvents)
			if err != nil {
				return fmt.Errorf("error creating MQTT input: %w", err)
			}
			err = input.Subscribe()
			if err != nil {
				return fmt.Errorf("error subscribing MQTT input: %w", err)
			}

		default:
			return fmt.Errorf("unknown input type: %s", i.Type)
		}
	}

	return nil
}
