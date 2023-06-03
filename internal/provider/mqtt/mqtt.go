package mqtt

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/valensto/ostraka/internal/logger"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/workflow"
)

type MQTT struct {
	client    mqtt.Client
	connected chan bool
	params    workflow.MQTTParams
	name      string
}

func connect(name string, params workflow.MQTTParams) (*MQTT, error) {
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

	c := &MQTT{
		name:      name,
		params:    params,
		client:    client,
		connected: make(chan bool),
	}

	if params.KeepAlive {
		go c.keepalive()
	}

	return c, nil
}

func (c *MQTT) keepalive() {
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

			logger.Get().Info().Msgf("%s send mqtt keep-alive", c.name)
			if token := c.client.Publish("ping_topic", 0, false, "ping"); token.Wait() && token.Error() != nil {
				log.Warn().Msgf("%s mqtt keep-alive failed: %s", c.name, token.Error())
			}
			break
		}
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

func (c *MQTT) Disconnect() {
	c.client.Disconnect(500)
	close(c.connected)
}

func (c *MQTT) Connect() error {
	if c.client.IsConnected() {
		return nil
	}
	token := c.client.Connect()
	token.Wait()
	return token.Error()
}
