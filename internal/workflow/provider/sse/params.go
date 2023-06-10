package sse

import (
	"encoding/json"
	"fmt"
)

type Params struct {
	Endpoint string `json:"endpoint"`
}

func (sse *Params) validate() error {
	if sse.Endpoint == "" {
		return fmt.Errorf("sse endpoint is empty")
	}

	return nil
}

func unmarshalParams(bytes []byte) (*Params, error) {
	sse := Params{}
	err := json.Unmarshal(bytes, &sse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling params to type SSE got: %w ", err)
	}

	err = sse.validate()
	if err != nil {
		return nil, err
	}

	return &sse, nil
}
