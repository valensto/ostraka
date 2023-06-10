package mqtt

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/valensto/ostraka/internal/logger"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	client    mqtt.Client
	connected chan bool
	name      string
	params    *Params
}

func (m *MQTT) connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.params.Broker)
	opts.SetClientID(fmt.Sprintf("%s-%s", uuid.New(), "ostraka"))
	opts.SetUsername(m.params.User)
	opts.SetPassword(m.params.Password)
	opts.SetAutoReconnect(m.params.AutoReconnect)
	opts.SetDefaultPublishHandler(defaultPubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to mqtt broker: %w", token.Error())
	}

	m.client = client

	if m.params.KeepAlive {
		go m.keepalive()
	}

	return nil
}

func (m *MQTT) keepalive() {
	const period = 2 * time.Second
	var up, closed bool
	log := logger.Get()
	for {
		select {
		case up, closed = <-m.connected:
			if !closed {
				return
			}
			break
		case <-time.After(period):
			if !up { // skip until we are connected
				continue
			}

			logger.Get().Info().Msgf("%s send mqtt keep-alive", m.name)
			if token := m.client.Publish("ping_topic", 0, false, "ping"); token.Wait() && token.Error() != nil {
				log.Warn().Msgf("%s mqtt keep-alive failed: %s", m.name, token.Error())
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

func (m *MQTT) Disconnect() {
	m.client.Disconnect(500)
	close(m.connected)
}

func (m *MQTT) Connect() error {
	if m.client.IsConnected() {
		return nil
	}
	token := m.client.Connect()
	token.Wait()
	return token.Error()
}
