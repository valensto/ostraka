package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

func NewSubscriber(input workflow.Input) (*MQTT, error) {
	params, err := input.MQTTParams()
	if err != nil {
		return nil, err
	}

	c, err := connect(input.Name, params)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *MQTT) Subscribe(events chan<- map[string]any) error {
	token := c.client.Subscribe(c.params.Topic, 1, c.eventPubHandler(events))
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", c.params.Topic)
	}

	logger.Get().Info().Msgf("input %s of type MQTT registered. Listening from topic %s", c.name, c.params.Topic)
	return nil
}

func (c *MQTT) eventPubHandler(events chan<- map[string]any) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var data map[string]any
		err := json.Unmarshal(msg.Payload(), &data)
		if err != nil {
			logger.Get().Error().Msgf("error decoding message: %s", err)
			return
		}

		events <- data
		logger.Get().Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
	}
}
