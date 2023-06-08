package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/sse"
	"github.com/valensto/ostraka/internal/workflow"
)

type publisher interface {
	Publish(events <-chan workflow.Event) error
}

func (d dispatcher) registerOutputs() error {
	for _, output := range d.workflow.Outputs {
		d.outputs[output] = make(chan workflow.Event)

		p, err := d.getPublisher(output)
		if err != nil {
			return fmt.Errorf("error getting publisher: %w", err)
		}

		err = p.Publish(d.outputs[output])
		if err != nil {
			return fmt.Errorf("error registering publisher: %w", err)
		}
	}

	return nil
}

func (d dispatcher) getPublisher(output *workflow.Output) (publisher, error) {
	switch output.Destination {
	case workflow.SSE:
		return sse.NewPublisher(output, d.server)
	case workflow.MQTTPub:
		return mqtt.NewPublisher(output)
	default:
		return nil, fmt.Errorf("unknown output type: %s", output.Destination)
	}
}
