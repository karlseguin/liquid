package filters

import (
	"github.com/karlseguin/liquid/core"
)

// Creates an append filter
func AppendFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return Noop
	}
	return (&AppendFilter{parameters[0]}).Append
}

type AppendFilter struct {
	value core.Value
}

func (a *AppendFilter) Append(input interface{}, data map[string]interface{}) interface{} {
	var value string
	switch typed := a.value.Resolve(data).(type) {
	case string:
		value = typed
	default:
		value = string(core.ToBytes(typed))
	}
	switch typed := input.(type) {
	case string:
		return typed + value
	default:
		return core.ToString(input) + value
	}
}
