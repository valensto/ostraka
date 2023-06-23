package sse

import (
	"github.com/valensto/ostraka/internal/env"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/middleware"
)

func WebUIPublisher(config env.Webui, server *http.Server) (*Publisher, error) {
	p := &Publisher{
		server: server,
		params: &Params{
			Endpoint: "/webui/consumes",
		},
		authenticator: &middleware.Token{
			Token:      config.AuthToken,
			QueryParam: "token",
		},

		cors: &middleware.CORS{
			AllowedOrigins: config.AllowedOrigins,
			AllowedMethods: []string{"GET", "POST"},
		},
		clients:       make(map[client]struct{}),
		connecting:    make(chan client),
		disconnecting: make(chan client),
		bufSize:       2,
		eventCounter:  0,
	}

	endpoint := http.Endpoint{
		Method:  http.GET,
		Path:    p.params.Endpoint,
		Cors:    p.cors,
		Handler: p.endpoint(),
		Auth:    p.authenticator,
	}

	err := server.AddSubRouter(endpoint)
	if err != nil {
		return nil, err
	}

	p.listenConn()
	return p, nil
}
