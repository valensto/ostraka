package workflow

import (
	"encoding/json"
	"fmt"
)

type payload map[string]any

func (e payload) json() ([]byte, error) {
	b, ok := json.Marshal(e)
	if ok != nil {
		return nil, fmt.Errorf("error marshalling eventType to json: %w", ok)
	}
	return b, nil
}

type eventType struct {
	Format string  `json:"format" yaml:"format" validate:"required"`
	Fields []field `json:"fields" yaml:"fields" validate:"dive,required"`
}

type field struct {
	Name     string `json:"name" yaml:"name" validate:"required"`
	DataType string `json:"data_type" yaml:"data_type" validate:"required"`
	Required bool   `json:"required" yaml:"required"`
}
