package mqtt

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
)

type Subscriber struct {
	instance
	input *workflow.Input
}

func NewSubscriber(input *workflow.Input, params []byte) (*Subscriber, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	subscriber := Subscriber{
		instance: instance{
			name:   input.Name,
			params: p,
		},
		input: input,
	}

	err = subscriber.instance.connect()
	if err != nil {
		return nil, err
	}

	return &subscriber, nil
}

func (s *Subscriber) Input() *workflow.Input {
	return s.input
}

func (s *Subscriber) Subscribe(dispatch func(ctx context.Context, input *workflow.Input, data []byte) error, _ *server.Server) error {
	token := s.client.Subscribe(s.instance.params.Topic, 1, s.eventPubHandler(dispatch))
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", s.instance.params.Topic)
	}

	logger.Get().Info().Msgf("subscriber %s of type MQTT registered. Listening from topic %s", s.name, s.instance.params.Topic)
	return nil
}

func (s *Subscriber) eventPubHandler(dispatch func(ctx context.Context, input *workflow.Input, data []byte) error) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		logger.Get().Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())

		err := dispatch(context.Background(), s.input, msg.Payload())
		if err != nil {
			logger.Get().Error().Msgf("error dispatching message: %s", err)
			return
		}
	}
}
