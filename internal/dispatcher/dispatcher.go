package dispatcher

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/logger"
	"net/http"

	"github.com/valensto/ostraka/internal/workflow"
)

type extractor interface {
	Extract(_ context.Context) ([]*workflow.Workflow, error)
}

type dispatcher struct {
	workflow     *workflow.Workflow
	router       *chi.Mux
	inputEvents  chan map[string]any
	outputEvents map[string]chan []byte
}

func Dispatch(ctx context.Context, extractor extractor, port string) error {
	router := chi.NewRouter()
	workflows, err := extractor.Extract(ctx)
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		d := &dispatcher{
			workflow:     wf,
			router:       router,
			inputEvents:  make(chan map[string]any, len(wf.Inputs)),
			outputEvents: make(map[string]chan []byte),
		}

		go d.dispatchEvents()

		err := d.subscribeInputs()
		if err != nil {
			return err
		}

		err = d.registerOutputs()
		if err != nil {
			return err
		}
	}

	return http.ListenAndServe(":"+port, router)
}

func (d dispatcher) dispatchEvents() {
	log := logger.Get()
	for {
		select {
		case event := <-d.inputEvents:
			data, err := json.Marshal(event)
			if err != nil {
				log.Warn().Msgf("error marshaling event: %s", err.Error())
				continue
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
		}
	}
}
