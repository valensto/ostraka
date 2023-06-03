package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/sse"
	"github.com/valensto/ostraka/internal/workflow"
)

type Publisher interface {
	Register(events <-chan []byte) error
}

func (d dispatcher) registerOutputs() error {
	for _, output := range d.workflow.Outputs {
		d.outputEvents[output.Name] = make(chan []byte)

		publisher, err := d.getOutputProvider(output)
		if err != nil {
			return fmt.Errorf("error getting output publisher: %w", err)
		}

		err = publisher.Register(d.outputEvents[output.Name])
		if err != nil {
			return fmt.Errorf("error registering SSE output: %w", err)
		}
	}

	return nil
}

func (d dispatcher) getOutputProvider(output workflow.Output) (Publisher, error) {
	switch output.Destination {
	case workflow.SSE:
		return sse.New(output, d.server)
	case workflow.MQTTPub:
		return mqtt.NewPublisher(output)
	default:
		return nil, fmt.Errorf("unknown output type: %s", output.Destination)
	}
}
