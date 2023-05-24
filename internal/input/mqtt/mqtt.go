package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/valensto/ostraka/internal/config"
	"log"
)

type Input struct {
	client mqtt.Client
	params config.MQTTParams
	events chan<- map[string]any
}

func New(params config.MQTTParams, events chan<- map[string]any) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(params.Broker)
	opts.SetClientID(fmt.Sprintf("%s-%s", uuid.New(), "ostraka"))
	opts.SetUsername(params.User)
	opts.SetPassword(params.Password)

	opts.SetDefaultPublishHandler(defaultPubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to mqtt broker: %w", token.Error())
	}

	service := Input{
		client: client,
		params: params,
		events: events,
	}

	return service.subscribe()
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

	log.Printf("new mqtt input: %s registered", s.params.Topic)
	return nil
}

func (s *Input) eventPubHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		decoded, err := s.params.Decoder.Decode(msg.Payload())
		if err != nil {
			fmt.Printf("error decoding message: %v\n", err)
			return
		}

		data := map[string]any{}
		// check if the data is valid
		// map payload fields to the event config fields
		// use receiver on Decoder struct to add mappers logic
		// send mapped data

		s.events <- data
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
