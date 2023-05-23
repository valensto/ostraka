package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/valensto/ostraka/internal/config"
)

type Input struct {
	client mqtt.Client
	params config.MQTTParams
}

func New(params config.MQTTParams) (*Input, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(params.Broker)
	opts.SetClientID(fmt.Sprintf("%s-%s", params.User, "ostraka"))
	opts.SetUsername(params.User)
	opts.SetPassword(params.Password)

	opts.SetDefaultPublishHandler(defaultPubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error connecting to mqtt broker: %w", token.Error())
	}

	service := Input{
		client: client,
		params: params,
	}

	err := service.subscribe()
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *Input) Disconnect() {
	s.client.Disconnect(250)
}

func (s *Input) subscribe() error {
	token := s.client.Subscribe(s.params.Topic, 1, s.eventPubHandler())
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", s.params.Topic)
	}

	return nil
}

func (s *Input) eventPubHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		decoded, err := s.params.Decoder.Decode(msg.Payload())
		if err != nil {
			fmt.Printf("error decoding message: %v\n", err)
			return
		}

		fmt.Printf("Received message: %s from topic: %s\n", decoded, msg.Topic())
	}
}

var defaultPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}
