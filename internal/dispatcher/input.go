package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/input/mqtt"
	"github.com/valensto/ostraka/internal/input/webhook"
)

func (f file) proceedInputs() error {
	for _, i := range f.config.Inputs {
		switch i.Type {
		case config.Webhook:
			input, err := webhook.New(i, f.router, f.inputEvents)
			if err != nil {
				return fmt.Errorf("error creating webhook input: %w", err)
			}
			return input.Subscribe()

		case config.MQTT:
			input, err := mqtt.New(i, f.inputEvents)
			if err != nil {
				return fmt.Errorf("error creating MQTT input: %w", err)
			}
			return input.Subscribe()

		default:
			return fmt.Errorf("unknown input type: %s", i.Type)
		}
	}

	return nil
}
