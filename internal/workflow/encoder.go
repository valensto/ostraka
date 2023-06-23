package workflow

import (
	"fmt"
)

type encoder struct {
	Format Format `json:"format" yaml:"format" validate:"required"`
}

func (e encoder) encode(payload payload) ([]byte, error) {
	switch e.Format {
	case JSON:
		return payload.json()
	case HTML:
		return payload.html()
	default:
		return nil, fmt.Errorf("unknown encoder type: %s", e.Format)
	}
}
