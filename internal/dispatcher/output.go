package dispatcher

import (
	"fmt"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/output/sse"
)

func (f file) proceedOutputs() error {
	for _, output := range f.config.Outputs {
		switch output.Type {
		case config.SSE:
			params, err := output.ToSSEParams()
			if err != nil {
				return err
			}
			return sse.Register(params, f.router, f.outputEvents)

		default:
			return fmt.Errorf("unknown output type: %s", output.Type)
		}
	}

	return nil
}
