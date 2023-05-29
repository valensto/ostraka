package dispatcher

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/valensto/ostraka/internal/workflow"
	"github.com/valensto/ostraka/logger"
)

type dispatcher struct {
	workflow     workflow.Workflow
	router       *chi.Mux
	inputEvents  chan map[string]any
	outputEvents chan []byte
}

func Dispatch(workflows workflow.Workflows, port string) error {
	router := chi.NewRouter()

	for _, wf := range workflows {
		d := &dispatcher{
			workflow:     wf,
			router:       router,
			inputEvents:  make(chan map[string]any, len(wf.Inputs)),
			outputEvents: make(chan []byte, len(wf.Outputs)),
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

			log.Info().Msgf("event dispatched: %s", string(data))
			d.outputEvents <- data
		}
	}
}
