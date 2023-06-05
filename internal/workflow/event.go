package workflow

import "fmt"

type EventId string

func (id EventId) String() string {
	return string(id)
}

type Event struct {
	Id   EventId
	Data []byte
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
		return nil, fmt.Errorf("event type is empty")
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("event has no fields")
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
