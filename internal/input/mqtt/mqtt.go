package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/valensto/ostraka/internal/logger"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	"github.com/valensto/ostraka/internal/workflow"
)

type Client struct {
	client    mqtt.Client
	connected chan bool
	params    workflow.MQTTParams
	inputName string
	events    chan<- map[string]any
}

func New(input workflow.Input, events chan<- map[string]any) (*Client, error) {
	params, err := input.MQTTParams()
	if err != nil {
		return nil, err
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(params.Broker)
	opts.SetClientID(fmt.Sprintf("%s-%s", uuid.New(), "ostraka"))
	opts.SetUsername(params.User)
	opts.SetPassword(params.Password)
	opts.SetAutoReconnect(params.AutoReconnect)
	opts.SetDefaultPublishHandler(defaultPubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error connecting to mqtt broker: %w", token.Error())
	}

	c := &Client{
		inputName: input.Name,
		params:    params,
		events:    events,
		client:    client,
		connected: make(chan bool),
	}

	if params.KeepAlive {
		go c.keepalive()
	}

	return c, nil
}

func (c *Client) keepalive() {
	const period = 2 * time.Second
	var up, closed bool
	log := logger.Get()
	for {
		select {
		case up, closed = <-c.connected:
			if !closed {
				return
			}
			break
		case <-time.After(period):
			if !up { // skip until we are connected
				continue
			}

			logger.Get().Info().Msgf("%s send mqtt keep-alive", c.inputName)
			if token := c.client.Publish("ping_topic", 0, false, "ping"); token.Wait() && token.Error() != nil {
				log.Warn().Msgf("%s mqtt keep-alive failed: %s", c.inputName, token.Error())
			}
			break
		}
	}

}

func (c *Client) Subscribe() error {
	token := c.client.Subscribe(c.params.Topic, 1, c.eventPubHandler())
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", c.params.Topic)
	}

	logger.Get().Info().Msgf("input %s of type MQTT registered. Listening from topic %s", c.inputName, c.params.Topic)
	return nil
}

func (c *Client) eventPubHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var data map[string]any
		err := json.Unmarshal(msg.Payload(), &data)
		if err != nil {
			logger.Get().Error().Msgf("error decoding message: %s", err)
			return
		}

		c.events <- data
		logger.Get().Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
	}
}

var defaultPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.Get().Info().Msgf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logger.Get().Info().Msg("Connected to mqtt broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logger.Get().Warn().Msgf("Connection lost: %v", err)
}

func (c *Client) Disconnect() {
	c.client.Disconnect(500)
	close(c.connected)
}

func (c *Client) Connect() error {
	if c.client.IsConnected() {
		return nil
	}
	token := c.client.Connect()
	token.Wait()
	return token.Error()
}
