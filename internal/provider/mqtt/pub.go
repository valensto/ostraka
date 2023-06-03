package mqtt

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

type Pub struct {
	MQTT
	output workflow.Output
}

func NewPublisher(output workflow.Output) (*Pub, error) {
	params, err := output.MQTTParams()
	if err != nil {
		return nil, err
	}

	p := Pub{
		MQTT: MQTT{
			name:   output.Name,
			params: params,
		},
		output: output,
	}

	err = p.MQTT.connect()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Pub) Register(events <-chan []byte) error {
	l := logger.Get()
	l.Info().Msgf("output %s of type MQTT registered. Publishing to topic %s", p.name, p.params.Topic)

	go func() {
		for {
			select {
			case event := <-events:
				token := p.client.Publish(p.params.Topic, 1, false, event)
				token.Wait()
				if token.Error() != nil {
					l.Error().Msgf("error publishing to topic: %s", p.params.Topic)
				}
			}
		}
	}()

	return nil
}
