package workflow

import (
	"fmt"
	"github.com/valensto/ostraka/internal/provider"
)

type input struct {
	Name    string  `json:"name" yaml:"name" validate:"required"`
	Source  string  `json:"source" yaml:"source" validate:"required"`
	Decoder decoder `json:"decoder" yaml:"decoder" validate:"dive,required"`

	Params     any `json:"params" yaml:"params" validate:"required"`
	subscriber provider.Subscriber
	queue      chan []byte
}

func (i *input) loadSubscriber(opts provider.Options) error {
	subscriber, err := provider.NewSubscriber(i.Source, i.Params, opts)
	if err != nil {
		return err
	}

	i.subscriber = subscriber
	i.Params = nil
	return nil
}

func (wf *Workflow) loadInputs(opts provider.Options) error {
	for i, _ := range wf.Inputs {
		if err := wf.Inputs[i].loadSubscriber(opts); err != nil {
			return fmt.Errorf("error unmarshalling input %s got: %w", wf.Inputs[i].Name, err)
		}

		err := wf.Inputs[i].subscribe(wf.dispatch)
		if err != nil {
			return fmt.Errorf("error subscribing to input %s got: %w", wf.Inputs[i].Name, err)
		}

		wf.Inputs[i].Decoder.eventType = wf.EventType
	}

	return nil
}

func (i *input) subscribe(dispatch func(from *input, bytes []byte)) error {
	if i.subscriber == nil {
		return fmt.Errorf("input %s is not initialized", i.Name)
	}

	i.queue = make(chan []byte)

	go func() {
		for {
			select {
			case b := <-i.queue:
				dispatch(i, b)
			}
		}
	}()

	return i.subscriber.Subscribe(i.queue)
}
