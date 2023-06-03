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
	senders  map[*workflow.Output]chan []byte
}

func Dispatch(ctx context.Context, extractor extractor, port string) error {
	workflows, err := extractor.Extract(ctx)
	if err != nil {
		return err
	}

	s := server.New(port)

	for _, wf := range workflows {
		d := &dispatcher{
			workflow: wf,
			server:   s,
			senders:  make(map[*workflow.Output]chan []byte, len(wf.Outputs)),
		}

		err := d.subscribeInputs()
		if err != nil {
			return err
		}

		err = d.registerOutputs()
		if err != nil {
			return err
		}

		err = d.registerWebui()
		if err != nil {
			return err
		}
	}

	return s.Run(workflows)
}

func (d dispatcher) notifyWebUI(notifier server.Notifier, bytes []byte, err error) {
	// TODO: early return if webui is not enabled
	d.server.NotifyWebUI(d.workflow.Name, notifier, bytes, err)
}

func (d dispatcher) dispatch(from workflow.Input, data []byte) {
	var err error
	defer d.notifyWebUI(&from, data, err)

	event, err := from.Decoder.Decode(data)
	if err != nil {
		logger.Get().Error().Msgf("error decoding input %s: %s", from.Name, err)
		return
	}

	d.send(event)
}

func (d dispatcher) send(event map[string]any) {
	data, err := json.Marshal(event)
	if err != nil {
		logger.Get().Error().Msgf("error marshaling event: %s", err)
		return
	}

	for output, c := range d.senders {
		if !output.Condition.Match(event) {
			message := fmt.Errorf("event not matching output %s conditions", output.Name)
			d.notifyWebUI(output, data, message)
			logger.Get().Info().Msg(message.Error())
			continue
		}

		d.notifyWebUI(output, data, nil)
		c <- data
	}
}
