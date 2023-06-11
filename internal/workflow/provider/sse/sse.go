package sse

import (
	"bytes"
	"fmt"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow/middleware"
	"net/http"

	"github.com/valensto/ostraka/internal/workflow"
)

const SSE = "sse"

type Publisher struct {
	server *server.Server

	output        *workflow.Output
	params        *Params
	authenticator middleware.Authenticator
	cors          *middleware.CORS

	clients       map[client]bool
	connecting    chan client
	disconnecting chan client
	bufSize       uint
	eventCounter  uint32
}

type client chan []byte

func NewPublisher(output *workflow.Output, params []byte, middlewares *middleware.Middlewares) (*Publisher, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	publisher := Publisher{
		output:        output,
		params:        p,
		authenticator: nil,
		cors:          nil,

		clients:       make(map[client]bool),
		connecting:    make(chan client),
		disconnecting: make(chan client),
		bufSize:       2,
		eventCounter:  0,
	}

	if p.Auth != "" {
		publisher.authenticator, err = middlewares.Web.GetAuthenticator(p.Auth)
		if err != nil {
			return nil, err
		}
	}

	if p.CORS != "" {
		publisher.cors, err = middlewares.Web.GetCORS(p.CORS)
		if err != nil {
			return nil, err
		}
	}

	return &publisher, nil
}

func (o *Publisher) Output() *workflow.Output {
	return o.output
}

func (o *Publisher) Publish(events <-chan workflow.Event, mux *server.Server) error {
	if mux == nil {
		return fmt.Errorf("server is required to register publisher of type SSE")
	}

	o.server = mux
	endpoint := server.Endpoint{
		Method:  server.GET,
		Path:    o.params.Endpoint,
		Cors:    o.cors,
		Handler: o.endpoint(),
		Auth:    o.authenticator,
	}

	o.listen(events)

	err := o.server.AddSubRouter(endpoint)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("publisher %s of type SSE registered. Sending to endpoint %s", o.output.Name, o.params.Endpoint)
	return nil
}

func (o *Publisher) listen(events <-chan workflow.Event) {
	go func() {
		for {
			select {
			case cl := <-o.connecting:
				o.clients[cl] = true

			case cl := <-o.disconnecting:
				delete(o.clients, cl)

			case event := <-events:
				msg := format(fmt.Sprintf("%d", o.eventCounter), "message", event.Bytes())
				o.eventCounter++
				for cl := range o.clients {
					cl <- msg.Bytes()
				}
			}
		}
	}()
}

func (o *Publisher) endpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl, ok := w.(http.Flusher)
		if !ok {
			logger.Get().Error().Msg("error flushing response writer: flushing not supported")
			o.server.Respond(w, r, http.StatusNotImplemented, nil)
			return
		}

		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Headers", "Content-Type")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("Content-Type", "text/event-stream")

		cl := make(client, o.bufSize)
		o.connecting <- cl

		for {
			select {
			case <-r.Context().Done():
				o.disconnecting <- cl
				return

			case event := <-cl:
				_, _ = w.Write(event)
				fl.Flush()
			}
		}
	}
}

func format(id, event string, data []byte) *bytes.Buffer {
	var buffer bytes.Buffer

	if len(id) > 0 {
		buffer.WriteString(fmt.Sprintf("id: %s\n", id))
	}

	if len(event) > 0 {
		buffer.WriteString(fmt.Sprintf("event: %s\n", event))
	}

	if len(data) > 0 {
		buffer.WriteString(fmt.Sprintf("data: %s\n", string(data)))
	}

	buffer.WriteString("\n")

	return &buffer
}
