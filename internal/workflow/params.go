package workflow

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type parameter interface {
	validate() error
}

type WebhookParams struct {
	Endpoint string `json:"endpoint"`
}

func (w WebhookParams) validate() error {
	if w.Endpoint == "" {
		return fmt.Errorf("webhook endpoint is empty")
	}

	return nil
}

type MQTTParams struct {
	Broker        string `json:"broker"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Topic         string `json:"topic"`
	AutoReconnect bool   `json:"autoReconnect"`
	KeepAlive     bool   `json:"keepAlive"`
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

type SSEParams struct {
	Endpoint string `json:"endpoint"`
	Auth     Auth   `json:"auth,omitempty"`
}

func (sse SSEParams) validate() error {
	if sse.Endpoint == "" {
		return fmt.Errorf("sse endpoint is empty")
	}

	return nil
}

func unmarshalParams(marshalled []byte, params any) (err error) {
	t := reflect.TypeOf(params)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("type %T is not a pointer", params)
	}

	err = json.Unmarshal(marshalled, params)
	if err != nil {
		return fmt.Errorf("error unmarshalling params to type %T got: %w ", params, err)
	}

	return nil
}
