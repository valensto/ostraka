package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

type Subscriber struct {
	MQTT
	*workflow.Input
}

func NewSubscriber(input *workflow.Input) (*Subscriber, error) {
	params, err := input.MQTTParams()
	if err != nil {
		return nil, err
	}

	s := Subscriber{
		MQTT: MQTT{
			name:   input.Name,
			params: params,
		},
		Input: input,
	}

	err = s.MQTT.connect()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (m *Subscriber) Subscribe(dispatch func(from *workflow.Input, data []byte)) error {
	token := m.client.Subscribe(m.MQTT.params.Topic, 1, m.eventPubHandler(dispatch))
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", m.MQTT.params.Topic)
	}

	logger.Get().Info().Msgf("subscriber %s of type MQTT registered. Listening from topic %s", m.name, m.MQTT.params.Topic)
	return nil
}

func (m *Subscriber) eventPubHandler(dispatch func(from *workflow.Input, data []byte)) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		logger.Get().Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
		dispatch(m.Input, msg.Payload())
	}
}
