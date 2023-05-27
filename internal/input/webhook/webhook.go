package webhook

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/valensto/ostraka/internal/config"
)

type Input struct {
	router *chi.Mux
	config.Input
	params config.WebhookParams
	events chan<- map[string]any
}

func New(input config.Input, router *chi.Mux, events chan<- map[string]any) (*Input, error) {
	params, err := input.GetAsWebhookParams()
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

	log.Infof("new webhook input: %s registered with endpoint %s", i.Name, i.params.Endpoint)
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
			log.Errorf("error decoding webhook input: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		i.events <- decoded
		w.WriteHeader(http.StatusOK)
	}
}
