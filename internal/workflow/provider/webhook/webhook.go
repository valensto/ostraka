package webhook

import (
	"fmt"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"io"
	"net/http"

	"github.com/valensto/ostraka/internal/workflow"
)

type Subscriber struct {
	server *server.Server
	params *Params
	input  *workflow.Input
}

func NewSubscriber(input *workflow.Input, params []byte) (*Subscriber, error) {
	p, err := unmarshalWebhook(params)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		input:  input,
		params: p,
	}, nil
}

func (s *Subscriber) Input() *workflow.Input {
	return s.input
}

func (s *Subscriber) Subscribe(dispatch func(input *workflow.Input, data []byte) error, mux *server.Server) error {
	if mux == nil {
		return fmt.Errorf("server is required to register subscriber of type webhook")
	}

	s.server = mux
	endpoint := server.Endpoint{
		Method:  server.POST,
		Path:    s.params.Endpoint,
		Handler: s.endpoint(dispatch),
	}

	err := s.server.AddSubRouter(endpoint)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("subscriber %s of type webhook registered. Listening from endpoint %s", s.input.Name, s.params.Endpoint)
	return nil
}

func (s *Subscriber) endpoint(dispatch func(input *workflow.Input, data []byte) error) http.HandlerFunc {
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

		err = dispatch(s.input, bytes)
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
