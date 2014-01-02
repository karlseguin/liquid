package filters

import (
	"strconv"
	"time"
)

var defaultPlus = &IntPlusFilter{1, "1"}

// Creates a plus filter
func PlusFactory(parameters []string) Filter {
	if len(parameters) == 0 || parameters[0] == "1" {
		return defaultPlus.Plus
	}
	if i, err := strconv.Atoi(parameters[0]); err == nil {
		return (&IntPlusFilter{i, parameters[0]}).Plus
	}
	if f, err := strconv.ParseFloat(parameters[0], 64); err == nil {
		return (&FloatPlusFilter{f, parameters[0]}).Plus
	}
	return (&StringPlusFilter{parameters[0]}).Plus
}

// Plus filter for integer parameter
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
		return stringPlusInt(typed, p.plus, p.plusString)
	case []byte:
		return stringPlusInt(string(typed), p.plus, p.plusString)
	default:
		return input
	}
}

func stringPlusInt(s string, plus int, plusString string) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return i + plus
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f + float64(plus)
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

// Plus filter for float parameter
type FloatPlusFilter struct {
	plus       float64
	plusString string
}

func (p *FloatPlusFilter) Plus(input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return float64(typed) + p.plus
	case float64:
		return typed + p.plus
	case uint:
		return float64(typed) + p.plus
	case int64:
		return float64(typed) + p.plus
	case uint64:
		return float64(typed) + p.plus
	case string:
		return stringPlusFloat(typed, p.plus, p.plusString)
	case []byte:
		return stringPlusFloat(string(typed), p.plus, p.plusString)
	default:
		return input
	}
}

func stringPlusFloat(s string, plus float64, plusString string) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return float64(i) + plus
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f + plus
	}
	return s + plusString
}
