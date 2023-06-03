package workflow

import "fmt"

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

func (o *Output) SSEParams() (SSEParams, error) {
	params, ok := o.params.(SSEParams)
	if !ok {
		return SSEParams{}, fmt.Errorf("output params are not of type SSEParams")
	}

	return params, nil
}
