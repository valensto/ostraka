package workflow

import (
	"fmt"
)

type Output struct {
	Name        string
	Destination string
	Condition   *Condition
}

func UnmarshallOutput(name, destination string, condition *Condition) (*Output, error) {
	if name == "" {
		return nil, fmt.Errorf("output name is empty")
	}

	if destination == "" {
		return nil, fmt.Errorf("output destination is empty")
	}

	return &Output{
		Name:        name,
		Destination: destination,
		Condition:   condition,
	}, nil
}
