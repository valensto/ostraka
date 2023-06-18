package sse

import (
	"bytes"
	"fmt"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/middleware"
	stdHTTP "net/http"
)

const SSE = "sse"

type Publisher struct {
	server *http.Server

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

func NewPublisher(params []byte, s *http.Server, middlewares *middleware.Middlewares) (*Publisher, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	publisher := Publisher{
		server:        s,
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
		publisher.authenticator, err = middlewares.HTTP.Authenticator(p.Auth)
		if err != nil {
			return nil, err
		}
	}

	if p.CORS != "" {
		publisher.cors, err = middlewares.HTTP.Cors(p.CORS)
		if err != nil {
			return nil, err
		}
	}

	endpoint := http.Endpoint{
		Method:  http.GET,
		Path:    publisher.params.Endpoint,
		Cors:    publisher.cors,
		Handler: publisher.endpoint(),
		Auth:    publisher.authenticator,
	}

	err = s.AddSubRouter(endpoint)
	if err != nil {
		return nil, err
	}

	publisher.listenConn()
	logger.Get().Info().Msgf("publisher of type SSE registered. Sending to endpoint %s", publisher.params.Endpoint)
	return &publisher, nil
}

func (p *Publisher) Provider() string {
	return SSE
}

func (p *Publisher) Publish(b []byte) {
	msg := format(fmt.Sprintf("%d", p.eventCounter), "message", b)
	p.eventCounter++
	for cl := range p.clients {
		cl <- msg.Bytes()
	}

	logger.Get().Info().Msgf("event published to endpoint %s", p.params.Endpoint)
}

func (p *Publisher) listenConn() {
	go func() {
		for {
			select {
			case cl := <-p.connecting:
				p.clients[cl] = true

			case cl := <-p.disconnecting:
				delete(p.clients, cl)
			}
		}
	}()
}

func (p *Publisher) endpoint() stdHTTP.HandlerFunc {
	return func(w stdHTTP.ResponseWriter, r *stdHTTP.Request) {
		fl, ok := w.(stdHTTP.Flusher)
		if !ok {
			logger.Get().Error().Msg("error flushing response writer: flushing not supported")
			p.server.Respond(w, r, stdHTTP.StatusNotImplemented, nil)
			return
		}

		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Headers", "Content-Type")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("Content-Type", "text/event-stream")

		cl := make(client, p.bufSize)
		p.connecting <- cl

		for {
			select {
			case <-r.Context().Done():
				p.disconnecting <- cl
				return

			case e := <-cl:
				_, _ = w.Write(e)
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
