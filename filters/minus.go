package filters

import (
	"github.com/karlseguin/liquid/core"
)

var defaultMinus = (&IntPlusFilter{-1}).Plus

// Creates a minus filter
func MinusFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return defaultMinus
	}
	switch typed := parameters[0].(type) {
	case *core.StaticIntValue:
		return (&IntPlusFilter{-typed.Value}).Plus
	case *core.StaticFloatValue:
		return (&FloatPlusFilter{-typed.Value}).Plus
	case *core.DynamicValue:
		return (&DynamicMinusFilter{typed}).Minus
	default:
		if n, ok := core.ToInt(parameters[0].Underlying()); ok {
			return (&IntPlusFilter{-n}).Plus
		}
		return Noop
	}
}

type DynamicMinusFilter struct {
	value core.Value
}

func (p *DynamicMinusFilter) Minus(input interface{}, data map[string]interface{}) interface{} {
	resolved := p.value.Resolve(data)
	switch typed := resolved.(type) {
	case int:
		return plusInt(-typed, input)
	case float64:
		return plusFloat(-typed, input)
	default:
		return input
	}
}
