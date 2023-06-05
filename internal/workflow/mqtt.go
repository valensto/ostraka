package workflow

import "fmt"

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

func (o *Output) MQTTParams() (MQTTParams, error) {
	if o.Destination != MQTTPub {
		return MQTTParams{}, fmt.Errorf("output source is not MQTT")
	}

	return toMQTTParams(o.params)
}

func (i *Input) MQTTParams() (MQTTParams, error) {
	if i.Source != MQTTSub {
		return MQTTParams{}, fmt.Errorf("input source is not MQTT")
	}

	return toMQTTParams(i.params)
}

func toMQTTParams(params any) (MQTTParams, error) {
	mqtt, ok := params.(MQTTParams)
	if !ok {
		return MQTTParams{}, fmt.Errorf("input params are not of type MQTTParams")
	}

	return mqtt, nil
}
