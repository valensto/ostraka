package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/sse"
	"github.com/valensto/ostraka/internal/workflow"
)

type publisher interface {
	Publish(events <-chan []byte) error
}

func (d dispatcher) registerOutputs() error {
	for _, output := range d.workflow.Outputs {
		d.senders[&output] = make(chan []byte)

		p, err := d.getPublisher(output)
		if err != nil {
			return fmt.Errorf("error getting publisher: %w", err)
		}

		err = p.Publish(d.senders[&output])
		if err != nil {
			return fmt.Errorf("error registering publisher: %w", err)
		}
	}

	return nil
}

func (d dispatcher) registerWebui() error {
	output := workflow.WebUIOutput()

	p, err := sse.NewPublisher(output, d.server)
	if err != nil {
		return fmt.Errorf("error creating webui publisher: %w", err)
	}

	return p.Publish(d.server.Notifications())
}

func (d dispatcher) getPublisher(output workflow.Output) (publisher, error) {
	switch output.Destination {
	case workflow.SSE:
		return sse.NewPublisher(output, d.server)
	case workflow.MQTTPub:
		return mqtt.NewPublisher(output)
	default:
		return nil, fmt.Errorf("unknown output type: %s", output.Destination)
	}
}
