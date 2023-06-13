package sse

import (
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/workflow/middleware"
)

func WebUIPublisher(config env.Webui) *Publisher {
	return &Publisher{
		params: &Params{
			Endpoint: "/webui/consume",
		},
		authenticator: &middleware.Token{
			Token:      config.AuthToken,
			QueryParam: "token",
		},
		cors: &middleware.CORS{
			AllowedOrigins:   config.AllowedOrigins,
			AllowedMethods:   []string{"GET"},
			AllowedHeaders:   nil,
			AllowCredentials: false,
			MaxAge:           3000,
		},
		clients:       make(map[client]bool),
		connecting:    make(chan client),
		disconnecting: make(chan client),
		bufSize:       2,
		eventCounter:  0,
	}
}
