package workflow

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	SSE = "sse"
)

type Event struct {
	Type   string  `yaml:"type" validate:"required"`
	Fields []Field `yaml:"fields" validate:"required,dive,required"`
}

type Output struct {
	Name       string      `yaml:"name" validate:"required"`
	Type       string      `yaml:"type" validate:"required"`
	Params     interface{} `yaml:"params" validate:"required"`
	Conditions []Condition `yaml:"conditions"`
}

type SSEParams struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	Auth     Auth   `yaml:"auth" validate:"omitempty"`
}

type Condition struct {
	Conditions []Condition `yaml:"conditions"`
	Source     string      `yaml:"source" validate:"required"`
	Field      string      `yaml:"field" validate:"required"`
	Operator   string      `yaml:"operator" validate:"required"`
	Value      string      `yaml:"value" validate:"required"`
}

func (file *Workflow) populateOutputs() error {
	var parsedOutputs []Output

	for _, output := range file.Outputs {
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

	file.Outputs = parsedOutputs
	return nil
}

func (o Output) ToSSEParams() (SSEParams, error) {
	params, ok := o.Params.(SSEParams)
	if !ok {
		return SSEParams{}, fmt.Errorf("output params are not of type SSEParams")
	}

	return params, nil
}
