package workflow

type Condition struct {
	Field      string      `yaml:"field,omitempty"`
	Operator   string      `yaml:"operator"`
	Value      any         `yaml:"value,omitempty"`
	Conditions []Condition `yaml:"conditions,omitempty"`
}

func (c *Condition) Match(source map[string]any) bool {
	if c == nil {
		return true
	}

	if len(c.Conditions) > 0 {
		switch c.Operator {
		case "and":
			for _, subCondition := range c.Conditions {
				if !subCondition.Match(source) {
					return false
				}
			}
			return true

		case "or":
			for _, subCondition := range c.Conditions {
				if subCondition.Match(source) {
					return true
				}
			}
			return false

		default:
			return false
		}
	}

	return c.matchOperator(source)
}

func (c *Condition) matchOperator(source map[string]any) bool {
	v, ok := source[c.Field]
	if !ok {
		return false
	}

	switch c.Operator {
	case "eq":
		return v == c.Value
	case "ne":
		return v != c.Value
	case "gt", "lt", "gte", "lte":
		return compareNumbers(v, c.Value, c.Operator)
	case "in":
		return containsValue(c.Value, v)
	case "nin":
		return !containsValue(c.Value, v)
	case "exists":
		return v != nil
	case "nexists":
		return v == nil
	default:
		return false
	}
}

func compareNumbers(a, b any, operator string) bool {
	na, ok := a.(int)
	if !ok {
		return false
	}

	nb, ok := b.(int)
	if !ok {
		return false
	}

	switch operator {
	case "gt":
		return na > nb
	case "lt":
		return na < nb
	case "gte":
		return na >= nb
	case "lte":
		return na <= nb
	default:
		return false
	}
}

func containsValue(values any, value any) bool {
	arr, ok := values.([]interface{})
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
