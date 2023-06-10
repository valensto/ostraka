package workflow

type Input struct {
	Name    string
	Source  Source
	Decoder Decoder
}

func UnmarshallInput(name, source string, decoder Decoder, event *EventType) (*Input, error) {
	src, err := getSource(source)
	if err != nil {
		return nil, err
	}

	i := &Input{
		Name:    name,
		Source:  src,
		Decoder: decoder,
	}

	i.Decoder.event = event
	return i, nil
}
