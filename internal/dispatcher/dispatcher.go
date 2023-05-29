package dispatcher

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/valensto/ostraka/internal/workflow"
)

type file struct {
	config       workflow.Workflow
	router       *chi.Mux
	inputEvents  chan map[string]any
	outputEvents chan []byte
}

func newFile(conf workflow.Workflow, router *chi.Mux) *file {
	return &file{
		config:       conf,
		router:       router,
		inputEvents:  make(chan map[string]any, len(conf.Inputs)),
		outputEvents: make(chan []byte, len(conf.Outputs)),
	}
}

func Dispatch(workflows workflow.Workflows, port string) error {
	router := chi.NewRouter()

	for _, workflow := range workflows {
		f := newFile(workflow, router)

		go f.dispatchEvents()

		err := f.subscribeInputs()
		if err != nil {
			return err
		}

		err = f.registerOutputs()
		if err != nil {
			return err
		}
	}

	return http.ListenAndServe(":"+port, router)
}

func (f file) dispatchEvents() {
	for {
		select {
		case event := <-f.inputEvents:
			data, err := json.Marshal(event)
			if err != nil {
				log.Warnf("error marshaling event: %v", err)
				continue
			}

			log.Infof("event dispatched: %s", data)
			f.outputEvents <- data
		}
	}
}
