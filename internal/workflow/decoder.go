package workflow

import (
	"encoding/json"
	"fmt"
)

type Decoder struct {
	Format  string
	Mappers []Mapper
	event   *Event
}

type Mapper struct {
	Source string
	Target string
}

func (d Decoder) Decode(data []byte) (map[string]any, error) {
	if d.Format != "json" {
		return nil, fmt.Errorf("unknown decoder type: %s", d.Format)
	}

	var source map[string]any
	err := json.Unmarshal(data, &source)
	if err != nil {
		return nil, fmt.Errorf("error decoding message: %w", err)
	}

	var decoded = map[string]any{}
	for _, field := range d.event.fields {
		sf, ok := d.findSourceField(field.name)
		if !ok && field.required {
			return nil, fmt.Errorf("required field %s not found", field.name)
		}

		v, ok := source[sf]
		if !ok && field.required {
			return nil, fmt.Errorf("required field %s not found", field.name)
		}

		decoded[field.name] = v
	}

	return decoded, nil
}

func (d Decoder) findSourceField(target string) (source string, found bool) {
	for _, mapper := range d.Mappers {
		if mapper.Target == target {
			return mapper.Source, true
		}
	}

	return "", false
}
