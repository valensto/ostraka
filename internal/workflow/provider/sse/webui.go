package sse

import (
	"github.com/valensto/ostraka/internal/workflow"
)

func WebUIPublisher() *Publisher {
	return &Publisher{
		output: &workflow.Output{
			Name:        "webui",
			Destination: SSE,
		},
		params: &Params{
			Endpoint: "/webui/consume",
		},
		clients:       make(map[client]bool),
		connecting:    make(chan client),
		disconnecting: make(chan client),
		bufSize:       2,
		eventCounter:  0,
	}
}
