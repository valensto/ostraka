package workflow

import "fmt"

type Input struct {
	Name    string
	Source  string
	Decoder Decoder
}

func UnmarshallInput(name, source string, decoder Decoder, event *EventType) (*Input, error) {
	if name == "" {
		return nil, fmt.Errorf("input name is empty")
	}

	if source == "" {
		return nil, fmt.Errorf("input source is empty")
	}

	i := &Input{
		Name:    name,
		Source:  source,
		Decoder: decoder,
	}

	i.Decoder.event = event
	return i, nil
}
