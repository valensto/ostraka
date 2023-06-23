package workflow

import (
	"fmt"
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

func (o operator) String() string {
	return string(o)
}

func (o operator) isValid() bool {
	switch o {
	case And, Or, Gt, Lt, Eq, Ne, Gte, Lte, In, Nin:
		return true
	default:
		return false
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

type condition struct {
	Field      string      `json:"field,omitempty"`
	Operator   operator    `json:"operator"`
	Value      any         `json:"value,omitempty"`
	Conditions []condition `json:"conditions,omitempty"`
}

func (c *condition) newFromChildren(field string, operator operator, value any, conditions ...condition) (*condition, error) {
	if len(conditions) > 0 && (!operator.isTopOperator() || !operator.isValid()) {
		return nil, fmt.Errorf("invalid operator %s", operator)
	}

	if len(conditions) > 0 {
		return &condition{
			Operator:   operator,
			Conditions: conditions,
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

	return &condition{
		Field:      field,
		Operator:   operator,
		Value:      value,
		Conditions: conditions,
	}, nil
}

func (c *condition) computeConditions() (*condition, error) {
	updatedCs := make([]*condition, len(c.Conditions))
	for i, cs := range c.Conditions {
		nc, err := cs.computeConditions()
		if err != nil {
			return nil, fmt.Errorf("error converting condition yaml: %w", err)
		}
		updatedCs[i] = nc
	}

	return c.newFromChildren(c.Field, c.Operator, c.Value, c.Conditions...)
}

func (c *condition) match(event payload) bool {
	if c == nil {
		return true
	}

	if len(c.Conditions) > 0 {
		switch c.Operator {
		case And:
			for _, subCondition := range c.Conditions {
				if !subCondition.match(event) {
					return false
				}
			}
			return true

		case Or:
			for _, subCondition := range c.Conditions {
				if subCondition.match(event) {
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

func (c *condition) matchOperator(event payload) bool {
	v, ok := event[c.Field]
	if !ok {
		return false
	}

	switch c.Operator {
	case Eq:
		return v == c.Value
	case Ne:
		return v != c.Value
	case Gt, Lt, Gte, Lte:
		return compareNumbers(v, c.Value, c.Operator)
	case In:
		return containsValue(c.Value, v)
	case Nin:
		return !containsValue(c.Value, v)
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
