package sse

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/valensto/ostraka/internal/workflow"
	"github.com/valensto/ostraka/logger"
)

type sse struct {
	router        *chi.Mux
	params        workflow.SSEParams
	clients       map[client]bool
	connecting    chan client
	disconnecting chan client
	bufSize       uint
	eventCounter  uint32
}

type client chan []byte

func Register(output workflow.Output, router *chi.Mux, events <-chan []byte) error {
	params, err := output.ToSSEParams()
	if err != nil {
		return err
	}

	sse := sse{
		router:        router,
		params:        params,
		clients:       make(map[client]bool),
		connecting:    make(chan client),
		disconnecting: make(chan client),
		bufSize:       2,
		eventCounter:  0,
	}

	sse.listen(events)
	sse.router.Get(params.Endpoint, sse.endpoint())

	logger.Get().Info().Msgf("output %s of type SSE registered. Sending to endpoint %s", output.Name, params.Endpoint)
	return nil
}

func (s sse) endpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl, ok := w.(http.Flusher)
		if !ok {
			logger.Get().Error().Msg("error flushing response writer: flushing not supported")
			http.Error(w, "Flushing not supported", http.StatusNotImplemented)
			return
		}

		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Headers", "Content-Type")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("Content-Type", "text/event-stream")

		cl := make(client, s.bufSize)
		s.connecting <- cl

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				s.disconnecting <- cl
				return

			case event := <-cl:
				_, _ = w.Write(event)
				fl.Flush()
			}
		}
	}
}

func (s sse) listen(events <-chan []byte) {
	go func() {
		for {
			select {
			case cl := <-s.connecting:
				s.clients[cl] = true

			case cl := <-s.disconnecting:
				delete(s.clients, cl)

			case event := <-events:
				msg := format(fmt.Sprintf("%d", s.eventCounter), "message", event)
				s.eventCounter++
				for cl := range s.clients {
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
