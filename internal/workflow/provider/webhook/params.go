package webhook

import (
	"encoding/json"
	"fmt"
)

type Params struct {
	Endpoint string `json:"endpoint"`
}

func (w *Params) validate() error {
	if w.Endpoint == "" {
		return fmt.Errorf("webhook endpoint is empty")
	}

	return nil
}

func unmarshalWebhook(bytes []byte) (*Params, error) {
	w := Params{}
	err := json.Unmarshal(bytes, &w)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling params to type Webhook got: %w ", err)
	}

	err = w.validate()
	if err != nil {
		return nil, err
	}

	return &w, nil
}
