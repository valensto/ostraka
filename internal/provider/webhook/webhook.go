package webhook

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"io"
	"net/http"

	"github.com/valensto/ostraka/internal/workflow"
)

type Subscriber struct {
	server *server.Server
	params workflow.WebhookParams
	*workflow.Input
}

func NewSubscriber(input *workflow.Input, server *server.Server) (*Subscriber, error) {
	params, err := input.WebhookParams()
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		server: server,
		Input:  input,
		params: params,
	}, nil
}

func (i *Subscriber) Subscribe(dispatch func(input *workflow.Input, data []byte) error) error {
	endpoint := server.Endpoint{
		Method:  server.POST,
		Path:    i.params.Endpoint,
		Handler: i.endpoint(dispatch),
	}

	err := i.server.AddSubRouter(endpoint)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("subscriber %s of type webhook registered. Listening from endpoint %s", i.Name, i.params.Endpoint)
	return nil
}

func (i *Subscriber) endpoint(dispatch func(input *workflow.Input, data []byte) error) http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			i.server.Respond(w, r, http.StatusBadRequest, response{
				Message: "error reading request body",
			})
			return
		}

		err = dispatch(i.Input, bytes)
		if err != nil {
			i.server.Respond(w, r, http.StatusBadRequest, response{
				Message: "error dispatching event",
			})
			return
		}

		i.server.Respond(w, r, http.StatusOK, response{
			Message: "event dispatched",
		})
	}
}
