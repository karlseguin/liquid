package filters

import (
	"strconv"
)

// Creates a time filter
func TimesFactory(parameters []string) Filter {
	if len(parameters) == 0 {
		return Noop
	}
	if i, err := strconv.Atoi(parameters[0]); err == nil {
		return (&IntTimesFilter{i}).Times
	}
	if f, err := strconv.ParseFloat(parameters[0], 64); err == nil {
		return (&FloatTimesFilter{f}).Times
	}
	return Noop
}

type IntTimesFilter struct {
	times int
}

// Multiples two numbers
func (t *IntTimesFilter) Times(input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return typed * t.times
	case float64:
		return typed * float64(t.times)
	case uint:
		return typed * uint(t.times)
	case int64:
		return typed * int64(t.times)
	case uint64:
		return typed * uint64(t.times)
	case string:
		return stringTimesInt(typed, t.times)
	case []byte:
		return stringTimesInt(string(typed), t.times)
	default:
		return input
	}
}

func stringTimesInt(s string, times int) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return i * times
	}
	return s
}

type FloatTimesFilter struct {
	times float64
}

// Multiples two numbers
func (t *FloatTimesFilter) Times(input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return float64(typed) * t.times
	case float64:
		return typed * t.times
	case uint:
		return float64(typed) * t.times
	case int64:
		return float64(typed) * t.times
	case uint64:
		return float64(typed) * t.times
	case string:
		return stringTimesFloat(typed, t.times)
	case []byte:
		return stringTimesFloat(string(typed), t.times)
	default:
		return input
	}
}

func stringTimesFloat(s string, times float64) interface{} {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f * times
	}
	if i, err := strconv.Atoi(s); err == nil {
		return float64(i) * times
	}
	return s
}
