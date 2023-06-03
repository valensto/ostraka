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
	workflow     *workflow.Workflow
	server       *server.Server
	outputEvents map[string]chan []byte
	globalEvents chan []globalEvent
}

type globalEvent struct {
	WorkflowName string `json:"workflow_name"`
	SourceType   string `json:"source_type"`
	SourceName   string `json:"source_name"`
	Payload      []byte `json:"payload"`
	State        string `json:"state"`
}

func Start(ctx context.Context, extractor extractor, port string) error {
	s := server.New(port)

	workflows, err := extractor.Extract(ctx)
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		d := &dispatcher{
			workflow:     wf,
			server:       s,
			outputEvents: make(map[string]chan []byte, len(wf.Outputs)),
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

	return s.Serve()
}

func (d dispatcher) dispatch(from workflow.Input, bytes []byte) {
	log := logger.Get()

	event, err := from.Decoder.Decode(bytes)
	if err != nil {
		log.Error().Msgf("error decoding event: %s", err.Error())
		return
	}

	err = d.send(event)
	if err != nil {
		log.Error().Msgf("error sending event: %s", err.Error())
		return
	}
}

func (d dispatcher) send(event map[string]any) error {
	log := logger.Get()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	for outputName, c := range d.outputEvents {
		output, ok := d.workflow.Outputs[outputName]
		if !ok {
			log.Warn().Msgf("output %s not found", outputName)
			continue
		}

		match := output.Condition.Match(event)
		if !match {
			continue
		}

		log.Info().Msgf("event dispatched: %s to %s", data, outputName)
		c <- data
	}

	return nil
}
