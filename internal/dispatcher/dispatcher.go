package dispatcher

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/input/mqtt"
	"github.com/valensto/ostraka/internal/input/webhook"
	"github.com/valensto/ostraka/internal/output/sse"
)

type Dispatcher struct {
	router *chi.Mux
	conf   config.Config
}

func New(conf config.Config, router *chi.Mux) *Dispatcher {
	return &Dispatcher{
		conf:   conf,
		router: router,
	}
}

func (d Dispatcher) Dispatch() error {
	for _, file := range d.conf {
		err := d.proceedFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d Dispatcher) proceedFile(file config.File) error {
	events := make(chan map[string]any, len(file.Inputs))

	err := d.proceedInputs(file.Inputs, events)
	if err != nil {
		return err
	}

	err = d.proceedOutputs(file.Outputs, events)
	if err != nil {
		return err
	}

	return nil
}

func (d Dispatcher) proceedOutputs(outputs []config.Output, events <-chan map[string]any) error {
	for _, output := range outputs {
		switch output.Type {
		case "sse":
			params, err := output.ToSSEParams()
			if err != nil {
				return err
			}
			return sse.New(params, d.router, events)

		default:
			return fmt.Errorf("unknown output type: %s", output.Type)
		}
	}

	return nil
}

func (d Dispatcher) proceedInputs(inputs []config.Input, events chan<- map[string]any) error {
	for _, input := range inputs {
		switch input.Type {
		case "webhook":
			params, err := input.ToWebhookParams()
			if err != nil {
				return err
			}
			return webhook.New(params, d.router, events)

		case "mqtt":
			params, err := input.ToMQTTParams()
			if err != nil {
				return err
			}

			err = mqtt.New(params, events)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown input type: %s", input.Type)
		}
	}

	return nil
}
