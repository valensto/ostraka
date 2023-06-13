package webhook

import (
	"context"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
	"github.com/valensto/ostraka/internal/workflow/middleware"
	"io"
	"net/http"
)

const Webhook = "webhook"

type Subscriber struct {
	server *server.Server
	params *Params

	authenticator middleware.Authenticator
	cors          *middleware.CORS
}

func NewSubscriber(params []byte, server *server.Server, middlewares *middleware.Middlewares) (*Subscriber, error) {
	p, err := unmarshalWebhook(params)
	if err != nil {
		return nil, err
	}

	s := Subscriber{
		server:        server,
		params:        p,
		authenticator: nil,
		cors:          nil,
	}

	if p.Auth != "" {
		s.authenticator, err = middlewares.HTTP.Authenticator(p.Auth)
		if err != nil {
			return nil, err
		}
	}

	if p.CORS != "" {
		s.cors, err = middlewares.HTTP.Cors(p.CORS)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (s *Subscriber) Subscribe(dispatch func(ctx context.Context, input *workflow.Input, data []byte) error) error {
	endpoint := server.Endpoint{
		Method:  server.POST,
		Path:    s.params.Endpoint,
		Cors:    s.cors,
		Handler: s.endpoint(dispatch),
		Auth:    s.authenticator,
	}

	err := s.server.AddSubRouter(endpoint)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("subscriber of type webhook registered. Listening from endpoint %s", s.params.Endpoint)
	return nil
}

func (s *Subscriber) endpoint(dispatch func(ctx context.Context, input *workflow.Input, data []byte) error) http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			s.server.Respond(w, r, http.StatusBadRequest, response{
				Message: "error reading request body",
			})
			return
		}

		// TODO: add input
		err = dispatch(r.Context(), nil, bytes)
		if err != nil {
			s.server.Respond(w, r, http.StatusBadRequest, response{
				Message: "error dispatching event",
			})
			return
		}

		s.server.Respond(w, r, http.StatusOK, response{
			Message: "event dispatched",
		})
	}
}
