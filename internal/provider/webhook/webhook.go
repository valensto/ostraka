package webhook

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"io"
	"net/http"

	"github.com/valensto/ostraka/internal/workflow"
)

type Input struct {
	server *server.Server
	params workflow.WebhookParams
	workflow.Input
}

func New(input workflow.Input, server *server.Server) (*Input, error) {
	params, err := input.WebhookParams()
	if err != nil {
		return nil, err
	}

	return &Input{
		server: server,
		Input:  input,
		params: params,
	}, nil
}

func (i *Input) Subscribe(dispatch func(from workflow.Input, bytes []byte)) error {
	endpoint := server.Endpoint{
		Method:  server.POST,
		Path:    i.params.Endpoint,
		Handler: i.endpoint(dispatch),
	}

	err := i.server.AddSubRouter(endpoint)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("input %s of type webhook registered. Listening from endpoint %s", i.Name, i.params.Endpoint)
	return nil
}

func (i *Input) endpoint(dispatch func(from workflow.Input, bytes []byte)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dispatch(i.Input, bytes)
		w.WriteHeader(http.StatusOK)
	}
}
