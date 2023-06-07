package dispatcher

import (
	"encoding/json"
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
	outputs   map[*workflow.Output]chan []byte
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
			outputs:   make(map[*workflow.Output]chan []byte, len(wf.Outputs)),
			collector: collector.New(wf, consumer),
		}

		err := d.registerInputs()
		if err != nil {
			return err
		}

		err = d.registerOutputs()
		if err != nil {
			return err
		}
	}

	return s.Run()
}

func (d dispatcher) dispatch(input *workflow.Input, data []byte) {
	collect := d.collector.Collect(input, data)

	event, err := input.Decoder.Decode(data)
	if err != nil {
		collect.WithError(fmt.Errorf("error decoding input: %w", err)).Send()
		return
	}

	marshalled, err := json.Marshal(event)
	if err != nil {
		collect.WithError(fmt.Errorf("error marshalling event: %w", err)).Send()
		return
	}

	for output, c := range d.outputs {
		if !output.Condition.Match(event) {
			collect.
				WithError(fmt.Errorf("event does not match output %s condition", output.Name)).
				WithLogLevel(zerolog.InfoLevel).Send()
			continue
		}

		collect.WithOutput(output, marshalled).Send()
		c <- marshalled
	}
}
