package filters

import (
	"github.com/karlseguin/liquid/core"
	"strconv"
)

// Creates a time filter
func TimesFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return Noop
	}
	switch typed := parameters[0].(type) {
	case *core.StaticIntValue:
		return (&IntTimesFilter{typed.Value}).Times
	case *core.StaticFloatValue:
		return (&FloatTimesFilter{typed.Value}).Times
	case *core.DynamicValue:
		return (&DynamicTimesFilter{typed}).Times
	default:
		if n, ok := core.ToInt(parameters[0].Underlying()); ok {
			return (&IntTimesFilter{n}).Times
		}
		return Noop
	}
}

type DynamicTimesFilter struct {
	value core.Value
}

func (t *DynamicTimesFilter) Times(input interface{}, data map[string]interface{}) interface{} {
	resolved := t.value.Resolve(data)
	switch typed := resolved.(type) {
	case int:
		return timesInt(typed, input)
	case float64:
		return timesFloat(typed, input)
	default:
		return input
	}
}

type IntTimesFilter struct {
	times int
}

// Multiples two numbers
func (t *IntTimesFilter) Times(input interface{}, data map[string]interface{}) interface{} {
	return timesInt(t.times, input)
}

func timesInt(times int, input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return typed * times
	case float64:
		return typed * float64(times)
	case uint:
		return typed * uint(times)
	case int64:
		return typed * int64(times)
	case uint64:
		return typed * uint64(times)
	case string:
		return stringTimesInt(typed, times)
	case []byte:
		return stringTimesInt(string(typed), times)
	default:
		return input
	}
}

func stringTimesInt(s string, times int) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return i * times
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f * float64(times)
	}
	return s
}

type FloatTimesFilter struct {
	times float64
}

// Multiples two numbers
func (t *FloatTimesFilter) Times(input interface{}, data map[string]interface{}) interface{} {
	return timesFloat(t.times, input)
}

func timesFloat(times float64, input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return float64(typed) * times
	case float64:
		return typed * times
	case uint:
		return float64(typed) * times
	case int64:
		return float64(typed) * times
	case uint64:
		return float64(typed) * times
	case string:
		return stringTimesFloat(typed, times)
	case []byte:
		return stringTimesFloat(string(typed), times)
	default:
		return input
	}
}

func stringTimesFloat(s string, times float64) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return float64(i) * times
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f * times
	}
	return s
}
