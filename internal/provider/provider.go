package provider

import (
	"encoding/json"
	"fmt"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/middleware"
	"github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/smtp"
	"github.com/valensto/ostraka/internal/provider/sse"
	"github.com/valensto/ostraka/internal/provider/webhook"
)

type Subscriber interface {
	Subscribe(event chan<- []byte) error
}

type Publisher interface {
	Publish(b []byte)
}

type Options struct {
	middlewares *middleware.Middlewares
	server      *http.Server
}

func NewOptions(server *http.Server, middlewares *middleware.Middlewares) (Options, error) {
	if err := middlewares.LoadAuthenticators(); err != nil {
		return Options{}, fmt.Errorf("error loading authenticators: %w", err)
	}

	return Options{
		middlewares: middlewares,
		server:      server,
	}, nil
}

func NewSubscriber(src string, params any, opts Options) (Subscriber, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling input params: %w", err)
	}

	switch src {
	case webhook.Webhook:
		return webhook.NewSubscriber(b, opts.server, opts.middlewares)

	case mqtt.MQTT:
		return mqtt.NewSubscriber(b)

	default:
		return nil, fmt.Errorf("unknown subscriber type: %s", src)
	}
}

func NewPublisher(dst string, params any, opts Options) (Publisher, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling input params: %w", err)
	}

	switch dst {
	case sse.SSE:
		return sse.NewPublisher(b, opts.server, opts.middlewares)

	case mqtt.MQTT:
		return mqtt.NewPublisher(b)

	case smtp.SMTP:
		return smtp.NewPublisher(b)

	default:
		return nil, fmt.Errorf("unknown publisher type: %s", dst)
	}
}
