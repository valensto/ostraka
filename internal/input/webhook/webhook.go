package webhook

import (
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/logger"
	"io"
	"net/http"

	"github.com/valensto/ostraka/internal/workflow"
)

type Input struct {
	router *chi.Mux
	workflow.Input
	params workflow.WebhookParams
	events chan<- map[string]any
}

func New(input workflow.Input, router *chi.Mux, events chan<- map[string]any) (*Input, error) {
	params, err := input.WebhookParams()
	if err != nil {
		return nil, err
	}

	return &Input{
		router: router,
		Input:  input,
		params: params,
		events: events,
	}, nil
}

func (i *Input) Subscribe() error {
	i.router.Post(i.params.Endpoint, i.endpoint())

	logger.Get().Info().Msgf("input %s of type webhook registered. Listening from endpoint %s", i.Name, i.params.Endpoint)
	return nil
}

func (i *Input) endpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		decoded, err := i.Decoder.Decode(bytes)
		if err != nil {
			logger.Get().Error().Msgf("error decoding webhook input: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		i.events <- decoded
		w.WriteHeader(http.StatusOK)
	}
}
