package provider

import (
	"encoding/json"
	"fmt"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/middleware"
	"github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/sse"
	"github.com/valensto/ostraka/internal/provider/webhook"
)

type Subscriber interface {
	Subscribe(event chan<- []byte) error
	Provider() string
}

type Publisher interface {
	Publish(b []byte)
	Provider() string
}

type Options struct {
	Middlewares *middleware.Middlewares
	Server      *http.Server
}

func NewSubscriber(src string, params any, opts Options) (Subscriber, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling input params: %w", err)
	}

	switch src {
	case webhook.Webhook:
		return webhook.NewSubscriber(b, opts.Server, opts.Middlewares)

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
		return sse.NewPublisher(b, opts.Server, opts.Middlewares)

	case mqtt.MQTT:
		return mqtt.NewPublisher(b)

	default:
		return nil, fmt.Errorf("unknown publisher type: %s", dst)
	}
}
