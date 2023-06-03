package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

type Sub struct {
	MQTT
	input workflow.Input
}

func NewSubscriber(input workflow.Input) (*Sub, error) {
	params, err := input.MQTTParams()
	if err != nil {
		return nil, err
	}

	s := Sub{
		MQTT: MQTT{
			name:   input.Name,
			params: params,
		},
		input: input,
	}

	err = s.MQTT.connect()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (m *Sub) Subscribe(dispatch func(bytes []byte, from workflow.Input)) error {
	token := m.client.Subscribe(m.params.Topic, 1, m.eventPubHandler(dispatch))
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", m.params.Topic)
	}

	logger.Get().Info().Msgf("input %s of type MQTT registered. Listening from topic %s", m.name, m.params.Topic)
	return nil
}

func (m *Sub) eventPubHandler(dispatch func(bytes []byte, from workflow.Input)) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		dispatch(msg.Payload(), m.input)
		logger.Get().Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
	}
}
