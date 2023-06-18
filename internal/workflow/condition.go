package workflow

import (
	"fmt"
	"github.com/valensto/ostraka/internal/event"
)

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

type Condition struct {
	field      string
	operator   operator
	value      any
	conditions []*Condition
}

func NewCondition(field, operator string, value any, conditions ...*Condition) (*Condition, error) {
	op, err := getOperator(operator)
	if err != nil {
		return nil, err
	}

	if len(conditions) > 0 && !op.isTopOperator() {
		return nil, fmt.Errorf("invalid top operator %s", operator)
	}

	if len(conditions) > 0 {
		return &Condition{
			operator:   op,
			conditions: conditions,
		}, nil
	}

	if len(conditions) == 0 && field == "" {
		return nil, fmt.Errorf("invalid condition: field is empty")
	}

	if len(conditions) == 0 && operator == "" {
		return nil, fmt.Errorf("invalid condition: operator is empty")
	}

	if len(conditions) == 0 && value == nil {
		return nil, fmt.Errorf("invalid condition: value is empty")
	}

	return &Condition{
		field:      field,
		operator:   op,
		value:      value,
		conditions: conditions,
	}, nil
}

func (c *Condition) Match(event event.Payload) bool {
	if c == nil {
		return true
	}

	if len(c.conditions) > 0 {
		switch c.operator {
		case And:
			for _, subCondition := range c.conditions {
				if !subCondition.Match(event) {
					return false
				}
			}
			return true

		case Or:
			for _, subCondition := range c.conditions {
				if subCondition.Match(event) {
					return true
				}
			}
			return false

		default:
			return false
		}
	}

	return c.matchOperator(event)
}

func (c *Condition) matchOperator(event event.Payload) bool {
	v, ok := event[c.field]
	if !ok {
		return false
	}

	switch c.operator {
	case Eq:
		return v == c.value
	case Ne:
		return v != c.value
	case Gt, Lt, Gte, Lte:
		return compareNumbers(v, c.value, c.operator)
	case In:
		return containsValue(c.value, v)
	case Nin:
		return !containsValue(c.value, v)
	default:
		return false
	}
}

func compareNumbers(a, b any, operator operator) bool {
	na, ok := a.(int)
	if !ok {
		return false
	}

	nb, ok := b.(int)
	if !ok {
		return false
	}

	switch operator {
	case Gt:
		return na > nb
	case Lt:
		return na < nb
	case Gte:
		return na >= nb
	case Lte:
		return na <= nb
	default:
		return false
	}
}

func containsValue(values any, value any) bool {
	arr, ok := values.([]any)
	if !ok {
		return false
	}

	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}
