package workflow

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

type Source string

const (
	Webhook Source = "webhook"
	MQTT    Source = "mqtt"
)

func getSource(src string) (Source, error) {
	s := Source(src)
	if err := s.isValid(); err != nil {
		return "", err
	}
	return s, nil
}

func (src Source) String() string {
	return string(src)
}

func (src Source) isValid() error {
	switch src {
	case Webhook, MQTT:
		return nil
	default:
		return fmt.Errorf("invalid source: %s", src)
	}
}

type Destination string

const (
	SSE   Destination = "sse"
	Email Destination = "email"
)

func getDestination(dest string) (Destination, error) {
	d := Destination(dest)
	if err := d.isValid(); err != nil {
		return "", err
	}
	return d, nil
}

func (dest Destination) String() string {
	return string(dest)
}

func (dest Destination) isValid() error {
	switch dest {
	case SSE, Email:
		return nil
	default:
		return fmt.Errorf("invalid destination: %s", dest)
	}
}

type operator string

const (
	And operator = "and"
	Or  operator = "or"
	Gt  operator = "gt"
	Lt  operator = "lt"
	Eq  operator = "eq"
	Ne  operator = "ne"
	Gte operator = "gte"
	Lte operator = "lte"
	In  operator = "in"
	Nin operator = "nin"
)

func getOperator(op string) (operator, error) {
	o := operator(op)
	if err := o.isValid(); err != nil {
		return "", err
	}
	return o, nil
}

func (o operator) String() string {
	return string(o)
}

func (o operator) isValid() error {
	switch o {
	case And, Or, Gt, Lt, Eq, Ne, Gte, Lte, In, Nin:
		return nil
	default:
		return fmt.Errorf("invalid operator: %s", o)
	}
}

func (o operator) isTopOperator() bool {
	switch o {
	case And, Or:
		return true
	default:
		return false
	}
}
