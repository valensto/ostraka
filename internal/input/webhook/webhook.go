package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config"
	"log"
	"net/http"
)

type Input struct {
	params config.WebhookParams
	router *chi.Mux
}

func New(params config.WebhookParams, router *chi.Mux, events chan<- map[string]any) error {
	i := Input{
		params: params,
		router: router,
	}

	i.router.Post(params.Endpoint, i.endpoint(events))

	log.Printf("new webhook input: %s registered", params.Endpoint)
	return nil
}

func (i *Input) endpoint(events chan<- map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := i.decode(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events <- data
		w.WriteHeader(http.StatusOK)
	}
}

func (i *Input) decode(_ http.ResponseWriter, r *http.Request) (map[string]any, error) {
	if r.ContentLength == 0 {
		return nil, nil
	}

	var data map[string]any
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error decoding request body")
	}

	// check if the data is valid
	// map payload fields to the event config fields
	// use receiver on Decoder struct to add mappers logic
	// return mapped data

	return data, nil
}
