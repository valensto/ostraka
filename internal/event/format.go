package event

import "fmt"

type Format string

const (
	JSON Format = "json"
)

func getFormat(f string) (Format, error) {
	s := Format(f)
	if err := s.isValid(); err != nil {
		return "", err
	}
	return s, nil
}

func (f Format) String() string {
	return string(f)
}

func (f Format) isValid() error {
	switch f {
	case JSON:
		return nil
	default:
		return fmt.Errorf("invalid format: %s", f)
	}
}
