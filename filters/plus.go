package filters

import (
	"github.com/karlseguin/liquid/core"
	"strconv"
	"time"
)

var defaultPlus = (&IntPlusFilter{1}).Plus

// Creates a plus filter
func PlusFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return defaultPlus
	}
	switch typed := parameters[0].(type) {
	case *core.StaticIntValue:
		return (&IntPlusFilter{typed.Value}).Plus
	case *core.StaticFloatValue:
		return (&FloatPlusFilter{typed.Value}).Plus
	case *core.DynamicValue:
		return (&DynamicPlusFilter{typed}).Plus
	default:
		if n, ok := core.ToInt(parameters[0].Underlying()); ok {
			return (&IntPlusFilter{n}).Plus
		}
		return Noop
	}
}

type DynamicPlusFilter struct {
	value core.Value
}

func (p *DynamicPlusFilter) Plus(input interface{}, data map[string]interface{}) interface{} {
	resolved := p.value.Resolve(data)
	switch typed := resolved.(type) {
	case int:
		return plusInt(typed, input)
	case float64:
		return plusFloat(typed, input)
	default:
		return input
	}
}

// Plus filter for integer parameter
type IntPlusFilter struct {
	plus int
}

func (p *IntPlusFilter) Plus(input interface{}, data map[string]interface{}) interface{} {
	return plusInt(p.plus, input)
}

func plusInt(plus int, input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return typed + plus
	case float64:
		return typed + float64(plus)
	case uint:
		return typed + uint(plus)
	case int64:
		return typed + int64(plus)
	case uint64:
		return typed + uint64(plus)
	case time.Time:
		return typed.Add(time.Minute * time.Duration(plus))
	case string:
		return stringPlusInt(typed, plus)
	case []byte:
		return stringPlusInt(string(typed), plus)
	default:
		return input
	}
}

func stringPlusInt(s string, plus int) interface{} {
	if s == "now" || s == "today" {
		return core.Now().Add(time.Minute * time.Duration(plus))
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i + plus
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f + float64(plus)
	}
	return s
}

// Plus filter for float parameter
type FloatPlusFilter struct {
	plus float64
}

func (p *FloatPlusFilter) Plus(input interface{}, data map[string]interface{}) interface{} {
	return plusFloat(p.plus, input)
}

func plusFloat(plus float64, input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return float64(typed) + plus
	case float64:
		return typed + plus
	case uint:
		return float64(typed) + plus
	case int64:
		return float64(typed) + plus
	case uint64:
		return float64(typed) + plus
	case string:
		return stringPlusFloat(typed, plus)
	case []byte:
		return stringPlusFloat(string(typed), plus)
	default:
		return input
	}
}

func stringPlusFloat(s string, plus float64) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return float64(i) + plus
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f + plus
	}
	return s
}
