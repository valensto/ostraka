package mqtt

import (
	"github.com/valensto/ostraka/internal/logger"
)

type Publisher struct {
	instance
}

func NewPublisher(params []byte) (*Publisher, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	publisher := Publisher{
		instance: instance{
			params: p,
		},
	}

	err = publisher.instance.connect()
	if err != nil {
		return nil, err
	}

	logger.Get().Info().Msgf("publisher of type MQTT registered. Publishing to topic %s", publisher.params.Topic)
	return &publisher, nil
}

func (p *Publisher) Provider() string {
	return MQTT
}

func (p *Publisher) Publish(b []byte) {
	l := logger.Get()
	token := p.client.Publish(p.instance.params.Topic, 1, false, b)
	token.Wait()
	if token.Error() != nil {
		l.Error().Msgf("error publishing to topic: %s", p.instance.params.Topic)
	}

	l.Info().Msgf("published message to topic: %s", p.instance.params.Topic)
}
