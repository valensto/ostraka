package event

import (
	"fmt"
)

type Encoder struct {
	format Format
}

func UnmarshalEncoder() (*Encoder, error) {
	return &Encoder{format: JSON}, nil
}

func (e Encoder) Encode(payload Payload) ([]byte, error) {
	switch e.format {
	case JSON:
		return payload.JSONEncode()
	default:
		return nil, fmt.Errorf("unknown encoder type: %s", e.format)
	}
}
