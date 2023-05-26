package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/output/sse"
)

func (f file) registerOutputs() error {
	for _, output := range f.config.Outputs {
		switch output.Type {
		case config.SSE:
			return sse.Register(output, f.router, f.outputEvents)

		default:
			return fmt.Errorf("unknown output type: %s", output.Type)
		}
	}

	return nil
}
