package dispatcher

import (
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
	outputs   map[*workflow.Output]chan workflow.Event
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
			outputs:   make(map[*workflow.Output]chan workflow.Event, len(wf.Publishers)),
			collector: collector.New(wf, consumer),
		}

		err := d.registerInputs(s)
		if err != nil {
			return err
		}

		err = d.registerOutputs(s)
		if err != nil {
			return err
		}
	}

	return s.Run()
}

func (d dispatcher) registerInputs(mux *server.Server) error {
	for _, s := range d.workflow.Subscribers {
		err := s.Subscribe(d.dispatch, mux)
		if err != nil {
			return fmt.Errorf("error registering subscriber: %w", err)
		}
	}

	return nil
}

func (d dispatcher) registerOutputs(mux *server.Server) error {
	for _, p := range d.workflow.Publishers {
		d.outputs[p.Output()] = make(chan workflow.Event)

		err := p.Publish(d.outputs[p.Output()], mux)
		if err != nil {
			return fmt.Errorf("error registering publisher: %w", err)
		}
	}

	return nil
}

func (d dispatcher) dispatch(input *workflow.Input, data []byte) error {
	collect := d.collector.Collect(input, data)

	event, err := input.Decoder.Decode(data)
	if err != nil {
		collect.WithError(fmt.Errorf("error decoding input: %w", err)).Send()
		return collect.Error()
	}

	for output, c := range d.outputs {
		if !output.Condition.Match(event) {
			collect.
				WithError(fmt.Errorf("event does not match output %s condition", output.Name)).
				WithLogLevel(zerolog.InfoLevel).Send()
			continue
		}

		collect.WithOutput(output, event).Send()
		c <- event
	}

	return nil
}
