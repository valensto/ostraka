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

type Condition struct {
	Field      string      `yaml:"field,omitempty"`
	Operator   string      `yaml:"operator"`
	Value      any         `yaml:"value,omitempty"`
	Conditions []Condition `yaml:"conditions,omitempty"`
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

func (c *Condition) IsMatching(source map[string]any) bool {
	if c == nil {
		return true
	}

	if len(c.Conditions) > 0 {
		switch c.Operator {
		case "and":
			for _, subCondition := range c.Conditions {
				if !subCondition.IsMatching(source) {
					return false
				}
			}
			return true
		case "or":
			for _, subCondition := range c.Conditions {
				if subCondition.IsMatching(source) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}

	return c.matchOperator(source)
}

func (c *Condition) matchOperator(source map[string]any) bool {
	v, ok := source[c.Field]
	if !ok {
		return false
	}

	switch c.Operator {
	case "eq":
		fmt.Printf("comparing %v with %v\n", v, c.Value)
		return v == c.Value
	case "ne":
		return v != c.Value
	case "gt", "lt", "gte", "lte":
		return compareNumbers(v, c.Value, c.Operator)
	case "in":
		return containsValue(c.Value, v)
	case "nin":
		return !containsValue(c.Value, v)
	case "exists":
		return v != nil
	case "nexists":
		return v == nil
	default:
		return false
	}
}

func compareNumbers(a, b any, operator string) bool {
	na, ok := a.(int)
	if !ok {
		return false
	}

	nb, ok := b.(int)
	if !ok {
		return false
	}

	switch operator {
	case "gt":
		return na > nb
	case "lt":
		return na < nb
	case "gte":
		return na >= nb
	case "lte":
		return na <= nb
	default:
		return false
	}
}

func containsValue(values any, value any) bool {
	arr, ok := values.([]interface{})
	if !ok {
		return false
	}

	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}
