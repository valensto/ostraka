package dispatcher

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config"
	"log"
)

type file struct {
	config       config.File
	router       *chi.Mux
	inputEvents  chan map[string]any
	outputEvents chan []byte
}

func newFile(conf config.File, router *chi.Mux) *file {
	return &file{
		config:       conf,
		router:       router,
		inputEvents:  make(chan map[string]any, len(conf.Inputs)),
		outputEvents: make(chan []byte, len(conf.Outputs)),
	}
}

func Dispatch(conf config.Config, router *chi.Mux) error {
	for _, file := range conf {
		f := newFile(file, router)

		go f.dispatchEvents()

		err := f.proceedInputs()
		if err != nil {
			return err
		}

		err = f.proceedOutputs()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f file) dispatchEvents() {
	for {
		select {
		case event := <-f.inputEvents:
			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("error marshaling event: %v", err)
				continue
			}

			f.outputEvents <- data
		}
	}
}
