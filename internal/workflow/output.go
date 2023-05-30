package workflow

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

const (
	SSE = "sse"
)

type Output struct {
	Name      string      `yaml:"name" validate:"required"`
	Type      string      `yaml:"type" validate:"required"`
	Params    interface{} `yaml:"params" validate:"required"`
	Condition *Condition  `yaml:"condition,omitempty"`
}

type SSEParams struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	Auth     Auth   `yaml:"auth" validate:"omitempty"`
}

func (wf *Workflow) setOutputs() error {
	var parsedOutputs []Output

	for _, output := range wf.Outputs {
		marshalled, err := yaml.Marshal(output.Params)
		if err != nil {
			return fmt.Errorf("error marshalling output params: %w", err)
		}

		switch output.Type {
		case SSE:
			var params SSEParams
			err := unmarshalParams(marshalled, &params)
			if err != nil {
				return err
			}
			output.Params = params

		default:
			return fmt.Errorf("unknown output type: %s", output.Type)
		}

		parsedOutputs = append(parsedOutputs, output)
	}

	wf.Outputs = parsedOutputs
	return nil
}

func (o Output) ToSSEParams() (SSEParams, error) {
	params, ok := o.Params.(SSEParams)
	if !ok {
		return SSEParams{}, fmt.Errorf("output params are not of type SSEParams")
	}

	return params, nil
}
