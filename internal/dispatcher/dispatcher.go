package dispatcher

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/valensto/ostraka/internal/collector"
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/consumer/webui"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
)

type dispatcher struct {
	workflow  *workflow.Workflow
	server    *server.Server
	collector *collector.Collector
}

func Dispatch(config *env.Config, workflows []*workflow.Workflow) error {
	s := server.New(config)
	consumer, err := webui.New(config.Webui, s, workflows)
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		d := &dispatcher{
			workflow:  wf,
			server:    s,
			collector: collector.New(wf, consumer),
		}

		d.registerInputs()
	}

	return s.Run()
}

func (d dispatcher) dispatch(ctx context.Context, input *workflow.Input, data []byte) error {
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
