package webhook

import (
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/middleware"
	"io"
	stdHTTP "net/http"
)

const Webhook = "webhook"

type Subscriber struct {
	server *http.Server
	params *Params

	authenticator middleware.Authenticator
	cors          *middleware.CORS
}

func NewSubscriber(params []byte, server *http.Server, middlewares *middleware.Middlewares) (*Subscriber, error) {
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

func (s *Subscriber) Subscribe(events chan<- []byte) error {
	endpoint := http.Endpoint{
		Method:  http.POST,
		Path:    s.params.Endpoint,
		Cors:    s.cors,
		Handler: s.endpoint(events),
		Auth:    s.authenticator,
	}

	err := s.server.AddSubRouter(endpoint)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("subscriber of type webhook registered. Listening from endpoint %s", s.params.Endpoint)
	return nil
}

func (s *Subscriber) endpoint(events chan<- []byte) stdHTTP.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}

	return func(w stdHTTP.ResponseWriter, r *stdHTTP.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			s.server.Respond(w, r, stdHTTP.StatusBadRequest, response{
				Message: "error reading request body",
			})
			return
		}

		events <- b
		s.server.Respond(w, r, stdHTTP.StatusOK, response{
			Message: "event dispatched",
		})
	}
}
