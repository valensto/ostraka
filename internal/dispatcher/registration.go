package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/workflow"
)

func (d dispatcher) registerInputs() error {
	for _, s := range d.workflow.Subscribers {
		err := s.Subscribe(d.dispatch, d.server)
		if err != nil {
			return fmt.Errorf("error registering subscriber: %w", err)
		}
	}

	return nil
}

func (d dispatcher) registerOutputs() error {
	for _, p := range d.workflow.Publishers {
		d.outputs[p.Output()] = make(chan workflow.Event)

		err := p.Publish(d.outputs[p.Output()], d.server)
		if err != nil {
			return fmt.Errorf("error registering publisher: %w", err)
		}
	}

	return nil
}
