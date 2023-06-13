package workflow

import (
	"encoding/json"
	"fmt"
)

type Event map[string]any

func (e Event) jsonEncode() ([]byte, error) {
	b, ok := json.Marshal(e)
	if ok != nil {
		return nil, fmt.Errorf("error marshalling eventType to json: %w", ok)
	}
	return b, nil
}

type EventType struct {
	format string
	fields []Field
}

type Field struct {
	name     string
	dataType string
	required bool
}

func UnmarshallEventType(format string, fields ...Field) (*EventType, error) {
	if format == "" {
		return nil, fmt.Errorf("eventType type is empty")
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("eventType has no fields")
	}

	return &EventType{
		format: format,
		fields: fields,
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
		name:     name,
		dataType: dataType,
		required: required,
	}, nil
}
