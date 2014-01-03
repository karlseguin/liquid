package filters

import (
	"github.com/karlseguin/liquid/core"
)

// Creates an append filter
func TruncateFactory(parameters []core.Value) Filter {
	if len(parameters) == 0 {
		return Noop
	}
	return (&TruncateFilter{parameters[0]}).Truncate
}

type TruncateFilter struct {
	value core.Value
}

func (t *TruncateFilter) Truncate(input interface{}, data map[string]interface{}) interface{} {
	length, ok := core.ToInt(t.value.Resolve(data))
	if ok == false {
		return input
	}

	var value string
	switch typed := input.(type) {
	case string:
		value = typed
	default:
		value = core.ToString(typed)
	}

	if length > len(value) {
		return input
	}
	return value[:length]
}
