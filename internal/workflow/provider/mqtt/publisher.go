package mqtt

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
)

type Publisher struct {
	instance
	output *workflow.Output
}

func NewPublisher(output *workflow.Output, params []byte) (*Publisher, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	publisher := Publisher{
		instance: instance{
			name:   output.Name,
			params: p,
		},
		output: output,
	}

	err = publisher.instance.connect()
	if err != nil {
		return nil, err
	}

	return &publisher, nil
}

func (p *Publisher) Output() *workflow.Output {
	return p.output
}

func (p *Publisher) Publish(events <-chan workflow.Event, _ *server.Server) error {
	l := logger.Get()
	l.Info().Msgf("publisher %s of type MQTT registered. Publishing to topic %s", p.name, p.instance.params.Topic)

	go func() {
		for {
			select {
			case event := <-events:
				token := p.client.Publish(p.instance.params.Topic, 1, false, event.Bytes())
				token.Wait()
				if token.Error() != nil {
					l.Error().Msgf("error publishing to topic: %s", p.instance.params.Topic)
				}
			}
		}
	}()

	return nil
}
