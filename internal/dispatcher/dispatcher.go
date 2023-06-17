package dispatcher

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/webui"
	workflow3 "github.com/valensto/ostraka/internal/workflow"
)

type dispatcher struct {
	workflow  *workflow3.Workflow
	server    *http.Server
	collector *workflow3.Collector
}

func Dispatch(config *env.Config, workflows []*workflow3.Workflow) error {
	s := http.New(config)
	consumer, err := webui.New(config.Webui, s, workflows)
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		d := &dispatcher{
			workflow:  wf,
			server:    s,
			collector: workflow3.New(wf, consumer),
		}

		d.registerInputs()
	}

	return s.Run()
}

func (d dispatcher) dispatch(ctx context.Context, input *workflow3.Input, data []byte) error {
	collect := d.collector.Collect(input, data)

	event, err := input.Decoder.Decode(ctx, data)
	if err != nil {
		collect.WithError(fmt.Errorf("error decoding input: %w", err)).Send()
		return collect.Error()
	}

	for _, output := range d.workflow.Outputs {
		err := output.Publish(event)
		collect.WithOutput(output, event).WithError(err).WithLogLevel(zerolog.InfoLevel).Send()
	}

	return nil
}
