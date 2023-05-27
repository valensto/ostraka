package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	Webhook = "webhook"
	MQTT    = "mqtt"
)

type Input struct {
	Name    string      `yaml:"name" validate:"required"`
	Type    string      `yaml:"type" validate:"required"`
	Decoder Decoder     `yaml:"decoder" validate:"required,dive,required"`
	Params  interface{} `yaml:"params" validate:"required"`
}

type Decoder struct {
	Type    string   `yaml:"type" validate:"required"`
	Mappers []Mapper `yaml:"mappers" validate:"required,dive,required"`
	event   Event    `yaml:"-"`
}

type Mapper struct {
	Source string `yaml:"source" validate:"required"`
	Target string `yaml:"target" validate:"required"`
}

type WebhookParams struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
}

type MQTTParams struct {
	Broker        string `yaml:"broker" validate:"required"`
	User          string `yaml:"user" validate:"required"`
	Password      string `yaml:"password" validate:"required"`
	Topic         string `yaml:"topic" validate:"required"`
	AutoReconnect bool   `yaml:"autoreconnect" validate:"required"`
	KeepAlive     bool   `yaml:"keepalive" validate:"required"`
}

func (file *Workflow) setInputs() error {
	var parsedInputs []Input

	for _, input := range file.Inputs {
		marshalled, err := yaml.Marshal(input.Params)
		if err != nil {
			return fmt.Errorf("error marshalling input params: %w", err)
		}

		switch input.Type {
		case Webhook:
			var params WebhookParams
			err := unmarshalParams(marshalled, &params)
			if err != nil {
				return err
			}
			input.Params = params

		case MQTT:
			var params MQTTParams
			err := unmarshalParams(marshalled, &params)
			if err != nil {
				return err
			}
			input.Params = params

		default:
			return fmt.Errorf("unknown input type: %s", input.Type)
		}

		input.Decoder.event = file.Event
		parsedInputs = append(parsedInputs, input)
	}

	file.Inputs = parsedInputs
	return nil
}

func (i Input) GetAsWebhookParams() (WebhookParams, error) {
	params, ok := i.Params.(WebhookParams)
	if !ok {
		return WebhookParams{}, fmt.Errorf("input params are not of type WebhookParams")
	}

	return params, nil
}

func (i Input) GetAsMQTTParams() (MQTTParams, error) {
	params, ok := i.Params.(MQTTParams)
	if !ok {
		return MQTTParams{}, fmt.Errorf("input params are not of type MQTTParams")
	}

	return params, nil
}

func (d Decoder) Decode(data []byte) (map[string]any, error) {
	if d.Type != "json" {
		return nil, fmt.Errorf("unknown decoder type: %s", d.Type)
	}

	var source map[string]any
	err := json.Unmarshal(data, &source)
	if err != nil {
		return nil, fmt.Errorf("error decoding message: %w", err)
	}

	var decoded = map[string]any{}
	for _, field := range d.event.Fields {
		sf, ok := d.findSourceField(field.Name)
		if !ok && field.Required {
			return nil, fmt.Errorf("required field %s not found", field.Name)
		}

		v, ok := source[sf]
		if !ok && field.Required {
			return nil, fmt.Errorf("required field %s not found", field.Name)
		}

		// TODO: check field type
		/*_, ok = value.(field.Type)
		if !ok {
			return nil, fmt.Errorf("field %s is not of type %s", field.Name, field.Type)
		}*/

		decoded[field.Name] = v
	}

	return decoded, nil
}

func (d Decoder) findSourceField(target string) (source string, found bool) {
	for _, mapper := range d.Mappers {
		if mapper.Target == target {
			return mapper.Source, true
		}
	}

	return "", false
}
