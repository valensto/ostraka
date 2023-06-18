package event

import (
	"encoding/json"
	"fmt"
)

type Payload map[string]any

func (e Payload) JSONEncode() ([]byte, error) {
	b, ok := json.Marshal(e)
	if ok != nil {
		return nil, fmt.Errorf("error marshalling eventType to json: %w", ok)
	}
	return b, nil
}

type Type struct {
	Format string
	Fields []Field
}

type Field struct {
	Name     string
	DataType string
	Required bool
}

func UnmarshallType(format string, fields ...Field) (*Type, error) {
	if format == "" {
		return nil, fmt.Errorf("eventType type is empty")
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("eventType has no fields")
	}

	return &Type{
		Format: format,
		Fields: fields,
	}, nil
}

func UnmarshallField(name, dataType string, required bool) (Field, error) {
	if name == "" {
		return Field{}, fmt.Errorf("field name is empty")
	}

	if dataType == "" {
		return Field{}, fmt.Errorf("field dataType is empty")
	}

	return Field{
		Name:     name,
		DataType: dataType,
		Required: required,
	}, nil
}
