package dispatcher

import (
	"fmt"

	"github.com/valensto/ostraka/internal/output/sse"
	"github.com/valensto/ostraka/internal/workflow"
)

func (f file) registerOutputs() error {
	for _, output := range f.config.Outputs {
		switch output.Type {
		case workflow.SSE:
			return sse.Register(output, f.router, f.outputEvents)

		default:
			return fmt.Errorf("unknown output type: %s", output.Type)
		}
	}

	return nil
}
