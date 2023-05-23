package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config"
	"log"
	"net/http"
)

type Output struct {
	router        *chi.Mux
	params        config.SSEParams
	clients       map[client]bool
	connecting    chan client
	disconnecting chan client
	bufSize       uint
	eventCounter  uint32
}

type client chan []byte

func New(params config.SSEParams, router *chi.Mux, events <-chan map[string]any) error {
	output := Output{
		router:        router,
		params:        params,
		clients:       make(map[client]bool),
		connecting:    make(chan client),
		disconnecting: make(chan client),
		bufSize:       2,
		eventCounter:  0,
	}

	output.listen(events)
	output.router.Get(params.Endpoint, output.endpoint())

	log.Printf("new sse output: %s registered", params.Endpoint)
	return nil
}

func (s Output) endpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl, ok := w.(http.Flusher)
		if !ok {
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

		cl := make(client, s.bufSize)
		s.connecting <- cl

		for {
			select {
			case <-done:
				s.disconnecting <- cl
				return

			case event := <-cl:
				_, _ = w.Write(event)
				fl.Flush()
			}
		}
	}
}

func (s Output) listen(events <-chan map[string]any) {
	go func() {
		for {
			select {
			case cl := <-s.connecting:
				s.clients[cl] = true

			case cl := <-s.disconnecting:
				delete(s.clients, cl)

			case event := <-events:
				data, err := json.Marshal(event)
				if err != nil {
					log.Printf("error marshaling event: %v", err)
					return
				}

				msg := format(fmt.Sprintf("%v", s.eventCounter), "message", data)
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
