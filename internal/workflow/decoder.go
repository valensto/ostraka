package workflow

import (
	"encoding/json"
	"fmt"
)

type decoder struct {
	Format    Format   `json:"format" yaml:"format" validate:"required"`
	Mappers   []mapper `json:"mappers" yaml:"mappers" validate:"required"`
	eventType eventType
}

type mapper struct {
	Source string `json:"source" yaml:"source" validate:"required"`
	Target string `json:"target" yaml:"target" validate:"required"`
}

func (d *decoder) decode(data []byte) (payload, error) {
	if d.Format != JSON {
		return nil, fmt.Errorf("unknown decoder type: %s", d.Format)
	}

	var sources map[string]any
	err := json.Unmarshal(data, &sources)
	if err != nil {
		return nil, fmt.Errorf("error decoding eventType: %w", err)
	}

	var e = make(payload, len(d.eventType.Fields))
	for _, f := range d.eventType.Fields {
		sf, ok := d.getSourceFieldByTarget(f.Name)
		if !ok && f.Required {
			return nil, fmt.Errorf("required field %s not found", f.Name)
		}

		v, ok := sources[sf]
		if !ok && f.Required {
			return nil, fmt.Errorf("required field %s not found", f.Name)
		}

		e[f.Name] = v
	}

	return e, nil
}

func (d *decoder) getSourceFieldByTarget(target string) (source string, found bool) {
	for _, m := range d.Mappers {
		if m.Target == target {
			return m.Source, true
		}
	}

	return "", false
}
