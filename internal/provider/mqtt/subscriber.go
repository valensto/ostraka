package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/logger"
)

type Subscriber struct {
	instance
}

func NewSubscriber(params []byte) (*Subscriber, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	subscriber := Subscriber{
		instance: instance{
			params: p,
		},
	}

	err = subscriber.instance.connect()
	if err != nil {
		return nil, err
	}

	return &subscriber, nil
}

func (s *Subscriber) Subscribe(events chan<- []byte) error {
	token := s.client.Subscribe(s.instance.params.Topic, 1, s.eventPubHandler(events))
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", s.instance.params.Topic)
	}

	logger.Get().Info().Msgf("subscriber of type MQTT registered. Listening from topic %s", s.instance.params.Topic)
	return nil
}

func (s *Subscriber) Provider() string {
	return MQTT
}

func (s *Subscriber) eventPubHandler(events chan<- []byte) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		logger.Get().Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
		events <- msg.Payload()
	}
}
