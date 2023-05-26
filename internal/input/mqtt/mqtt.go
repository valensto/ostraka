package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/valensto/ostraka/internal/config"
)

type Input struct {
	client mqtt.Client
	params config.MQTTParams
	events chan<- map[string]any
}

func New(input config.Input, events chan<- map[string]any) (*Input, error) {
	params, err := input.ToMQTTParams()
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("error connecting to mqtt broker: %w", token.Error())
	}

	return &Input{
		params: params,
		events: events,
		client: client,
	}, nil
}

func (i *Input) Subscribe() error {
	token := i.client.Subscribe(i.params.Topic, 1, i.eventPubHandler())
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %s", i.params.Topic)
	}

	log.Infof("new mqtt input: %s registered", i.params.Topic)
	return nil
}

func (i *Input) eventPubHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var data map[string]any
		err := json.Unmarshal(msg.Payload(), &data)
		if err != nil {
			log.Errorf("error decoding message: %s", err)
			return
		}

		i.events <- data
		log.Infof("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}
}

var defaultPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Infof("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info("Connected to mqtt broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Warningf("Connection lost: %v", err)
}
