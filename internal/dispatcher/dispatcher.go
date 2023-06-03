package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
)

type extractor interface {
	Extract(_ context.Context) ([]*workflow.Workflow, error)
}

type dispatcher struct {
	workflow *workflow.Workflow
	server   *server.Server
	events   map[string]chan []byte
}

func Run(ctx context.Context, extractor extractor, port string) error {
	workflows, err := extractor.Extract(ctx)
	if err != nil {
		return err
	}

	s := server.New(port)

	for _, wf := range workflows {
		d := &dispatcher{
			workflow: wf,
			server:   s,
			events:   make(map[string]chan []byte, len(wf.Outputs)),
		}

		err := d.subscribeInputs()
		if err != nil {
			return err
		}

		err = d.registerOutputs()
		if err != nil {
			return err
		}
	}

	return s.Serve(workflows)
}

func (d dispatcher) notifyWebUI(notifier server.Notifier, bytes []byte, err error) {
	d.server.NotifyWebUI(d.workflow.Name, notifier, bytes, err)
}

func (d dispatcher) dispatch(from workflow.Input, bytes []byte) {
	event, err := d.receive(from, bytes)
	if err != nil {
		return
	}

	d.send(event)
}

func (d dispatcher) receive(from workflow.Input, bytes []byte) (event map[string]any, err error) {
	defer d.notifyWebUI(&from, bytes, err)

	event, err = from.Decoder.Decode(bytes)
	if err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}

	return event, nil
}

func (d dispatcher) send(event map[string]any) {
	data, err := json.Marshal(event)
	if err != nil {
		logger.Get().Error().Msgf("error marshaling event: %s", err)
		return
	}

	for outputName, c := range d.events {
		output, ok := d.workflow.Outputs[outputName]
		if !ok {
			d.notifyWebUI(&output, data, fmt.Errorf("output %s not found", outputName))
			continue
		}

		match := output.Condition.Match(event)
		if !match {
			d.notifyWebUI(&output, data, fmt.Errorf("event not matching output %s conditions", outputName))
			continue
		}

		d.notifyWebUI(&output, data, nil)
		c <- data
	}
}
