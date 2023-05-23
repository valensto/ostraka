package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config"
	"net/http"
)

type Input struct {
	params  config.WebhookParams
	handler http.Handler
}

func New(params config.WebhookParams, router *chi.Mux) error {
	i := Input{
		params: params,
	}

	router.Post(params.Endpoint, i.Endpoint())
	return nil
}

func (i *Input) Endpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := i.Decode(w, r, &config.Decoder{}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (i *Input) Decode(_ http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.ContentLength == 0 {
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}

	return nil
}
