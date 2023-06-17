package event

import (
	"encoding/json"
	"fmt"
)

type Decoder struct {
	format    Format
	mappers   []Mapper
	eventType Type
}

type Mapper struct {
	Source string
	Target string
}

func UnmarshallDecoder(format string, mappers []Mapper, eventType Type) (*Decoder, error) {
	f, err := getFormat(format)
	if err != nil {
		return nil, err
	}

	return &Decoder{
		format:    f,
		mappers:   mappers,
		eventType: eventType,
	}, nil
}

func (d Decoder) Decode(data []byte) (Payload, error) {
	if d.format != JSON {
		return nil, fmt.Errorf("unknown decoder type: %s", d.format)
	}

	var source map[string]any
	err := json.Unmarshal(data, &source)
	if err != nil {
		return nil, fmt.Errorf("error decoding eventType: %w", err)
	}

	var e = make(Payload, len(d.eventType.Fields))
	for _, field := range d.eventType.Fields {
		sf, ok := d.getSourceFieldByTarget(field.Name)
		if !ok && field.Required {
			return nil, fmt.Errorf("required field %s not found", field.Name)
		}

		v, ok := source[sf]
		if !ok && field.Required {
			return nil, fmt.Errorf("required field %s not found", field.Name)
		}

		e[field.Name] = v
	}

	return e, nil
}

func (d Decoder) getSourceFieldByTarget(target string) (source string, found bool) {
	for _, mapper := range d.mappers {
		if mapper.Target == target {
			return mapper.Source, true
		}
	}

	return "", false
}
