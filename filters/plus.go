package filters

import (
	"strconv"
	"time"
)

var defaultPlus = &IntPlusFilter{1, "1"}

// Creates a plus filter
func PlusFactory(parameters []string) Filter {
	if len(parameters) == 0 {
		return defaultPlus.Plus
	}
	if i, err := strconv.Atoi(parameters[0]); err == nil {
		return (&IntPlusFilter{i, parameters[0]}).Plus
	}
	return (&StringPlusFilter{parameters[0]}).Plus
}

// Plus filter for numbers
type IntPlusFilter struct {
	plus       int
	plusString string
}

func (p *IntPlusFilter) Plus(input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return typed + p.plus
	case float64:
		return typed + float64(p.plus)
	case uint:
		return typed + uint(p.plus)
	case int64:
		return typed + int64(p.plus)
	case uint64:
		return typed + uint64(p.plus)
	case time.Time:
		return typed.Add(time.Minute * time.Duration(p.plus))
	case string:
		return plusString(typed, p.plus, p.plusString)
	case []byte:
		return plusString(string(typed), p.plus, p.plusString)
	default:
		return input
	}
}

func plusString(s string, plus int, plusString string) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return i + plus
	}
	return s + plusString
}

// plus filter for strings
type StringPlusFilter struct {
	plus string
}

func (p *StringPlusFilter) Plus(input interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return typed + p.plus
	case []byte:
		return string(typed) + p.plus
	default:
		return input
	}
}
