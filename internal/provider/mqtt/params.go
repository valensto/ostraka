package mqtt

import (
	"encoding/json"
	"fmt"
)

type Params struct {
	Broker        string `json:"broker"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Topic         string `json:"topic"`
	AutoReconnect bool   `json:"autoReconnect"`
	KeepAlive     bool   `json:"keepAlive"`
}

func (mqtt *Params) validate() error {
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

func unmarshalParams(bytes []byte) (*Params, error) {
	mqtt := Params{}
	err := json.Unmarshal(bytes, &mqtt)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling params to type MQTT got: %w ", err)
	}

	err = mqtt.validate()
	if err != nil {
		return nil, err
	}

	return &mqtt, nil
}
