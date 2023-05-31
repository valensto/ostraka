package workflow

import (
	"encoding/json"
	"fmt"
)

type Input struct {
	Name    string
	Source  Source
	Decoder Decoder
	params  any
}

func UnmarshallInput(name, source string, decoder Decoder, params any, event *Event) (*Input, error) {
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
