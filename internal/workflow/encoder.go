package workflow

import "fmt"

type Encoder struct {
	format Format
}

func (e Encoder) Encode(event Event) ([]byte, error) {
	switch e.format {
	case JSON:
		return event.jsonEncode()
	default:
		return nil, fmt.Errorf("unknown encoder type: %s", e.format)
	}
}
