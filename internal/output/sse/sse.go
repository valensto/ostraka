package sse

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/logger"
	"net/http"

	"github.com/valensto/ostraka/internal/workflow"
)

type Output struct {
	router        *chi.Mux
	name          string
	params        workflow.SSEParams
	clients       map[client]bool
	connecting    chan client
	disconnecting chan client
	bufSize       uint
	eventCounter  uint32
}

type client chan []byte

func New(output workflow.Output, router *chi.Mux, events <-chan []byte) (*Output, error) {
	params, err := output.SSEParams()
	if err != nil {
		return nil, err
	}

	o := &Output{
		router:        router,
		name:          output.Name,
		params:        params,
		clients:       make(map[client]bool),
		connecting:    make(chan client),
		disconnecting: make(chan client),
		bufSize:       2,
		eventCounter:  0,
	}
	o.listen(events)

	return o, nil
}

func (o *Output) Register() error {
	o.router.Get(o.params.Endpoint, o.endpoint())

	logger.Get().Info().Msgf("output %s of type SSE registered. Sending to endpoint %s", o.name, o.params.Endpoint)
	return nil
}

func (o *Output) endpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl, ok := w.(http.Flusher)
		if !ok {
			logger.Get().Error().Msg("error flushing response writer: flushing not supported")
			http.Error(w, "Flushing not supported", http.StatusNotImplemented)
			return
		}

		done := r.Context().Done()

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
			case <-done:
				o.disconnecting <- cl
				return

			case event := <-cl:
				_, _ = w.Write(event)
				fl.Flush()
			}
		}
	}
}

func (o *Output) listen(events <-chan []byte) {
	go func() {
		for {
			select {
			case cl := <-o.connecting:
				o.clients[cl] = true

			case cl := <-o.disconnecting:
				delete(o.clients, cl)

			case event := <-events:
				msg := format(fmt.Sprintf("%d", o.eventCounter), "message", event)
				o.eventCounter++
				for cl := range o.clients {
					cl <- msg.Bytes()
				}
			}
		}
	}()
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
