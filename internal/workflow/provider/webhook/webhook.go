package webhook

import (
	"context"
	"fmt"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
	"github.com/valensto/ostraka/internal/workflow/middleware"
	"io"
	"net/http"
)

const Webhook = "webhook"

type Subscriber struct {
	server        *server.Server
	params        *Params
	input         *workflow.Input
	authenticator middleware.Authenticator
	cors          *middleware.CORS
}

func NewSubscriber(input *workflow.Input, params []byte, middlewares *middleware.Middlewares) (*Subscriber, error) {
	p, err := unmarshalWebhook(params)
	if err != nil {
		return nil, err
	}

	s := Subscriber{
		input:         input,
		params:        p,
		authenticator: nil,
		cors:          nil,
	}

	if p.Auth != "" {
		s.authenticator, err = middlewares.Web.GetAuthenticator(p.Auth)
		if err != nil {
			return nil, err
		}
	}

	if p.CORS != "" {
		s.cors, err = middlewares.Web.GetCORS(p.CORS)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (s *Subscriber) Input() *workflow.Input {
	return s.input
}

func (s *Subscriber) Subscribe(dispatch func(ctx context.Context, input *workflow.Input, data []byte) error, mux *server.Server) error {
	if mux == nil {
		return fmt.Errorf("server is required to register subscriber of type webhook")
	}

	s.server = mux
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

	logger.Get().Info().Msgf("subscriber %s of type webhook registered. Listening from endpoint %s", s.input.Name, s.params.Endpoint)
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

		err = dispatch(r.Context(), s.input, bytes)
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
