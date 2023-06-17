package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt2 "github.com/valensto/ostraka/internal/provider/mqtt"
	"github.com/valensto/ostraka/internal/provider/webhook"
)

type Subscriber interface {
	Subscribe(dispatch func(ctx context.Context, input *Input, data []byte) error) error
}

type Input struct {
	Name    string
	Source  string
	Decoder *Decoder

	Subscriber Subscriber
}

func UnmarshallInput(name, source string, decoder *Decoder, params any, opts Options) (*Input, error) {
	if name == "" {
		return nil, fmt.Errorf("input name is empty")
	}

	subscriber, err := newSubscriber(source, params, opts)
	if err != nil {
		return nil, err
	}

	return &Input{
		Name:       name,
		Source:     source,
		Decoder:    decoder,
		Subscriber: subscriber,
	}, nil
}

func newSubscriber(src string, params any, opts Options) (Subscriber, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling input params: %w", err)
	}

	switch src {
	case webhook.Webhook:
		return webhook.NewSubscriber(b, opts.Server, opts.Middlewares)

	case mqtt2.MQTT:
		return mqtt2.NewSubscriber(b)

	default:
		return nil, fmt.Errorf("unknown subscriber type: %s", src)
	}
}

func (i *Input) Subscribe(dispatch func(ctx context.Context, input *Input, data []byte) error) {
	i.Subscriber.Subscribe(dispatch)
}
