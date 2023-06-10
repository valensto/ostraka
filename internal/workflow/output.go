package workflow

import (
	"fmt"
)

type Output struct {
	Name        string
	Destination Destination
	Condition   *Condition
}

func WebUIOutput() *Output {
	return &Output{
		Name:        "webui",
		Destination: SSE,
		/*params: params.SSE{
			Endpoint: "/webui/consume",
		},*/
	}
}

func UnmarshallOutput(name, destination string, condition *Condition) (*Output, error) {
	if name == "" {
		return nil, fmt.Errorf("output name is empty")
	}

	dest, err := getDestination(destination)
	if err != nil {
		return nil, err
	}

	return &Output{
		Name:        name,
		Destination: dest,
		Condition:   condition,
	}, nil
}
