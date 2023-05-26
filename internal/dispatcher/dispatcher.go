package dispatcher

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/valensto/ostraka/internal/config"
	"net/http"
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

func Dispatch(conf config.Config, port string) error {
	router := chi.NewRouter()

	for _, file := range conf {
		f := newFile(file, router)

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
