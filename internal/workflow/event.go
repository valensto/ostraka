package workflow

import (
	"encoding/json"
	"fmt"
	"strings"
)

type payload map[string]any

func (e payload) json() ([]byte, error) {
	b, ok := json.Marshal(e)
	if ok != nil {
		return nil, fmt.Errorf("error marshalling eventType to json: %w", ok)
	}
	return b, nil
}

func (e payload) html() ([]byte, error) {
	var builder strings.Builder
	builder.WriteString("<ul>\n")
	for key, value := range e {
		item := fmt.Sprintf("<li>%s: %v</li>\n", key, value)
		builder.WriteString(item)
	}
	builder.WriteString("</ul>")
	return []byte(builder.String()), nil
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
