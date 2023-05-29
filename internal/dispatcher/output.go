package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/output/sse"
	"github.com/valensto/ostraka/internal/workflow"
)

func (f file) registerOutputs() error {
	for _, output := range f.workflow.Outputs {
		switch output.Type {
		case workflow.SSE:
			err := sse.Register(output, f.router, f.outputEvents)
			if err != nil {
				return fmt.Errorf("error registering SSE output: %w", err)
			}
		default:
			return fmt.Errorf("unknown output type: %s", output.Type)
		}
	}

	return nil
}
