package mqtt

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

type Publisher struct {
	MQTT
	*workflow.Output
}

func NewPublisher(output *workflow.Output) (*Publisher, error) {
	params, err := output.MQTTParams()
	if err != nil {
		return nil, err
	}

	p := Publisher{
		MQTT: MQTT{
			name:   output.Name,
			params: params,
		},
		Output: output,
	}

	err = p.MQTT.connect()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Publisher) Publish(events <-chan workflow.Event) error {
	l := logger.Get()
	l.Info().Msgf("publisher %s of type MQTT registered. Publishing to topic %s", p.name, p.MQTT.params.Topic)

	go func() {
		for {
			select {
			case event := <-events:
				token := p.client.Publish(p.MQTT.params.Topic, 1, false, event.Bytes())
				token.Wait()
				if token.Error() != nil {
					l.Error().Msgf("error publishing to topic: %s", p.MQTT.params.Topic)
				}
			}
		}
	}()

	return nil
}
