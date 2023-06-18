package workflow

import (
	"fmt"
	"github.com/valensto/ostraka/internal/event"
	"github.com/valensto/ostraka/internal/provider"
)

type Input struct {
	Name    string
	Decoder *event.Decoder

	Subscriber provider.Subscriber
	queue      chan []byte
}

func UnmarshallInput(name, source string, decoder *event.Decoder, params any, opts provider.Options) (*Input, error) {
	if name == "" {
		return nil, fmt.Errorf("input name is empty")
	}

	subscriber, err := provider.NewSubscriber(source, params, opts)
	if err != nil {
		return nil, fmt.Errorf("error creating subscriber for output %s got: %w", name, err)
	}

	return &Input{
		Name:    name,
		Decoder: decoder,

		Subscriber: subscriber,
		queue:      make(chan []byte),
	}, nil
}

func (i *Input) listen(dispatch func(from *Input, bytes []byte)) error {
	go func() {
		for {
			select {
			case b := <-i.queue:
				dispatch(i, b)
			}
		}
	}()

	return i.Subscriber.Subscribe(i.queue)
}
