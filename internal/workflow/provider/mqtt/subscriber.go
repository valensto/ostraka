package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
)

type Subscriber struct {
	MQTT
	input *workflow.Input
}

func NewSubscriber(input *workflow.Input, params []byte) (*Subscriber, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	subscriber := Subscriber{
		MQTT: MQTT{
			name:   input.Name,
			params: p,
		},
		input: input,
	}

	err = subscriber.MQTT.connect()
	if err != nil {
		return nil, err
	}

	return &subscriber, nil
}

func (s *Subscriber) Input() *workflow.Input {
	return s.input
}

func (s *Subscriber) Subscribe(dispatch func(input *workflow.Input, data []byte) error, _ *server.Server) error {
	token := s.client.Subscribe(s.MQTT.params.Topic, 1, s.eventPubHandler(dispatch))
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", s.MQTT.params.Topic)
	}

	logger.Get().Info().Msgf("subscriber %s of type MQTT registered. Listening from topic %s", s.name, s.MQTT.params.Topic)
	return nil
}

func (s *Subscriber) eventPubHandler(dispatch func(input *workflow.Input, data []byte) error) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		logger.Get().Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
		err := dispatch(s.input, msg.Payload())
		if err != nil {
			logger.Get().Error().Msgf("error dispatching message: %s", err)
			return
		}
	}
}
