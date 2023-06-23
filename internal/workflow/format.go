package workflow

import "fmt"

type Format string

const (
	JSON Format = "json"
	HTML Format = "html"
)

func (f Format) String() string {
	return string(f)
}

func (f Format) isValid() error {
	switch f {
	case JSON:
		return nil
	case HTML:
		return nil
	default:
		return fmt.Errorf("invalid format: %s", f)
	}
}
