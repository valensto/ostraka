package workflow

import (
	"context"
	"encoding/json"
	"fmt"
)

type Decoder struct {
	format  Format
	mappers []Mapper
	event   *EventType
}

type Mapper struct {
	Source string
	Target string
}

func UnmarshallDecoder(format string, mappers []Mapper) (*Decoder, error) {
	f, err := getFormat(format)
	if err != nil {
		return nil, err
	}

	return &Decoder{
		format:  f,
		mappers: mappers,
	}, nil
}

func (d Decoder) Decode(_ context.Context, data []byte) (Event, error) {
	if d.format != JSON {
		return nil, fmt.Errorf("unknown decoder type: %s", d.format)
	}

	var source map[string]any
	err := json.Unmarshal(data, &source)
	if err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}

	var event = make(Event, len(d.event.fields))
	for _, field := range d.event.fields {
		sf, ok := d.getSourceFieldByTarget(field.name)
		if !ok && field.required {
			return nil, fmt.Errorf("required field %s not found", field.name)
		}

		v, ok := source[sf]
		if !ok && field.required {
			return nil, fmt.Errorf("required field %s not found", field.name)
		}

		event[field.name] = v
	}

	return event, nil
}

func (d Decoder) getSourceFieldByTarget(target string) (source string, found bool) {
	for _, mapper := range d.mappers {
		if mapper.Target == target {
			return mapper.Source, true
		}
	}

	return "", false
}
