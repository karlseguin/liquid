package filters

import (
	"github.com/karlseguin/liquid/core"
)

// Creates an prepend filter
func PrependFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return Noop
	}
	return (&PrependFilter{parameters[0]}).Prepend
}

type PrependFilter struct {
	value core.Value
}

func (a *PrependFilter) Prepend(input interface{}, data map[string]interface{}) interface{} {
	var value string
	switch typed := a.value.Resolve(data).(type) {
	case string:
		value = typed
	default:
		value = string(core.ToBytes(typed))
	}
	switch typed := input.(type) {
	case string:
		return value + typed
	default:
		return value + core.ToString(typed)
	}
}
