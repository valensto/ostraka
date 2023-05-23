package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

const webhookType = "webhook"
const mqttType = "mqtt"

type Input struct {
	Name   string      `yaml:"name" validate:"required"`
	Type   string      `yaml:"type" validate:"required"`
	Params interface{} `yaml:"params" validate:"required"`
}

type WebhookParams struct {
	Endpoint string  `yaml:"endpoint" validate:"required"`
	Decoder  Decoder `yaml:"decoder" validate:"required,dive,required"`
}

type MQTTParams struct {
	Broker   string  `yaml:"broker" validate:"required"`
	User     string  `yaml:"user" validate:"required"`
	Password string  `yaml:"password" validate:"required"`
	Topic    string  `yaml:"topic" validate:"required"`
	Decoder  Decoder `yaml:"decoder" validate:"required,dive,required"`
}

type Decoder struct {
	Type    string   `yaml:"type" validate:"required"`
	Mappers []Mapper `yaml:"mappers" validate:"required,dive,required"`
}

type Mapper struct {
	Source string `yaml:"source" validate:"required"`
	Target string `yaml:"target" validate:"required"`
}

func (file *File) populateInputs() error {
	var parsedInputs []Input

	for _, input := range file.Inputs {
		marshalled, err := yaml.Marshal(input.Params)
		if err != nil {
			return fmt.Errorf("error marshalling input params: %w", err)
		}

		switch input.Type {
		case webhookType:
			var params WebhookParams
			err := unmarshalParams(marshalled, &params)
			if err != nil {
				return err
			}
			input.Params = params

		case mqttType:
			var params MQTTParams
			err := unmarshalParams(marshalled, &params)
			if err != nil {
				return err
			}
			input.Params = params

		default:
			return fmt.Errorf("unknown input type: %s", input.Type)
		}

		parsedInputs = append(parsedInputs, input)
	}

	file.Inputs = parsedInputs
	return nil
}

func (i Input) ToWebhookParams() (WebhookParams, error) {
	params, ok := i.Params.(WebhookParams)
	if !ok {
		return WebhookParams{}, fmt.Errorf("input params are not of type WebhookParams")
	}

	return params, nil
}

func (i Input) ToMQTTParams() (MQTTParams, error) {
	params, ok := i.Params.(MQTTParams)
	if !ok {
		return MQTTParams{}, fmt.Errorf("input params are not of type MQTTParams")
	}

	return params, nil
}

func (d Decoder) Decode(data []byte) (Decoder, error) {
	if d.Type != "json" {
		return Decoder{}, fmt.Errorf("unknown decoder type: %s", d.Type)
	}

	err := json.Unmarshal(data, &d)
	if err != nil {
		return Decoder{}, err
	}

	return d, nil
}
