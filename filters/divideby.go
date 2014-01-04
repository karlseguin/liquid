package filters

import (
	"github.com/karlseguin/liquid/core"
)

func DivideByFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return Noop
	}
	switch typed := parameters[0].(type) {
	case *core.StaticIntValue:
		return (&FloatTimesFilter{1 / float64(typed.Value)}).Times
	case *core.StaticFloatValue:
		return (&FloatTimesFilter{1 / typed.Value}).Times
	case *core.DynamicValue:
		return (&DynamicDivideByFilter{typed}).DivideBy
	default:
		if n, ok := core.ToInt(parameters[0].Underlying()); ok {
			return (&FloatTimesFilter{1 / float64(n)}).Times
		}
		return Noop
	}
}

type DynamicDivideByFilter struct {
	value core.Value
}

func (t *DynamicDivideByFilter) DivideBy(input interface{}, data map[string]interface{}) interface{} {
	resolved := t.value.Resolve(data)
	switch typed := resolved.(type) {
	case int:
		return timesFloat(1/float64(typed), input)
	case float64:
		return timesFloat(1/typed, input)
	default:
		return input
	}
}
