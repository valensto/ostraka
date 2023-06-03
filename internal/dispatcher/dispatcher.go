package dispatcher

import (
	"encoding/json"
	"fmt"
	"github.com/valensto/ostraka/internal/collector"
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/consumer/webui"
	"github.com/valensto/ostraka/internal/logger"
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
			collector: collector.New(wf.Name, consumer),
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

	return s.Run()
}

func (d dispatcher) dispatch(from *workflow.Input, data []byte) {
	var err error
	defer d.collector.Collect(from, data, err)

	event, err := from.Decoder.Decode(data)
	if err != nil {
		err = fmt.Errorf("error decoding input %s: %s", from.Name, err)
		logger.Get().Error().Msg(err.Error())
		return
	}

	marshalled, err := json.Marshal(event)
	if err != nil {
		err = fmt.Errorf("error marshaling event: %s", err)
		logger.Get().Error().Msg(err.Error())
		return
	}

	for output, c := range d.outputs {
		if !output.Condition.Match(event) {
			err = fmt.Errorf("event not matching output %s conditions", output.Name)
			logger.Get().Info().Msg(err.Error())
			continue
		}

		d.collector.Collect(output, data, nil)
		c <- marshalled
	}
}
