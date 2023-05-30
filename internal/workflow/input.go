package workflow

import (
	"encoding/json"
	"fmt"
)

type Input struct {
	Name    string
	Source  Source
	Decoder Decoder
	params  interface{}
}

func UnmarshallInput(name, source string, decoder Decoder, params interface{}, event *Event) (*Input, error) {
	src, err := getSource(source)
	if err != nil {
		return nil, err
	}

	i := &Input{
		Name:    name,
		Source:  src,
		Decoder: decoder,
		params:  params,
	}

	err = i.unmarshallParams(event)
	if err != nil {
		return nil, err
	}

	return i, nil
}

type WebhookParams struct {
	Endpoint string `json:"endpoint"`
}

type MQTTParams struct {
	Broker        string `json:"broker"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Topic         string `json:"topic"`
	AutoReconnect bool   `json:"autoReconnect"`
	KeepAlive     bool   `json:"keepAlive"`
}

func (i *Input) unmarshallParams(e *Event) error {
	marshalled, err := json.Marshal(i.params)
	if err != nil {
		return fmt.Errorf("error marshalling input params: %w", err)
	}

	var params parameter
	switch i.Source {
	case Webhook:
		var wh WebhookParams
		err = unmarshalParams(marshalled, &wh)
		if err != nil {
			return err
		}

		params = wh
	case MQTT:
		var mqtt MQTTParams
		err = unmarshalParams(marshalled, &mqtt)
		if err != nil {
			return err
		}

		params = mqtt
	default:
		return fmt.Errorf("unknown input type: %s", i.Source)
	}

	i.params = params
	i.Decoder.event = e

	return params.validate()
}

func (i *Input) WebhookParams() (WebhookParams, error) {
	if i.Source != Webhook {
		return WebhookParams{}, fmt.Errorf("input source is not Webhook")
	}

	params, ok := i.params.(WebhookParams)
	if !ok {
		return WebhookParams{}, fmt.Errorf("input params are not of type WebhookParams")
	}

	return params, nil
}

func (w WebhookParams) validate() error {
	if w.Endpoint == "" {
		return fmt.Errorf("webhook endpoint is empty")
	}

	return nil
}

func (i *Input) MQTTParams() (MQTTParams, error) {
	if i.Source != MQTT {
		return MQTTParams{}, fmt.Errorf("input source is not MQTT")
	}

	params, ok := i.params.(MQTTParams)
	if !ok {
		return MQTTParams{}, fmt.Errorf("input params are not of type MQTTParams")
	}

	return params, nil
}

func (mqtt MQTTParams) validate() error {
	if mqtt.Broker == "" {
		return fmt.Errorf("mqtt broker is empty")
	}

	if mqtt.Topic == "" {
		return fmt.Errorf("mqtt topic is empty")
	}

	if mqtt.User == "" {
		return fmt.Errorf("mqtt user is empty")
	}

	if mqtt.Password == "" {
		return fmt.Errorf("mqtt password is empty")
	}

	return nil
}
