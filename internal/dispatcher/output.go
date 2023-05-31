package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/output/sse"
	"github.com/valensto/ostraka/internal/workflow"
)

type OutputProvider interface {
	Register() error
}

func (d dispatcher) registerOutputs() error {
	for _, output := range d.workflow.Outputs {
		d.outputEvents[output.Name] = make(chan []byte)

		provider, err := d.getOutputProvider(output)
		if err != nil {
			return fmt.Errorf("error getting output provider: %w", err)
		}

		err = provider.Register()
		if err != nil {
			return fmt.Errorf("error registering SSE output: %w", err)
		}
	}

	return nil
}

func (d dispatcher) getOutputProvider(o workflow.Output) (OutputProvider, error) {
	switch o.Destination {
	case workflow.SSE:
		return sse.New(o, d.router, d.outputEvents[o.Name])
	default:
		return nil, fmt.Errorf("unknown output type: %s", o.Destination)
	}
}
