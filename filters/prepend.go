package filters

import (
	"github.com/karlseguin/liquid/core"
)

// Creates an prepend filter
func PrependFactory(parameters []string) Filter {
	if len(parameters) == 0 || parameters[0] == "" {
		return Noop
	}
	return (&PrependFilter{parameters[0]}).Prepend
}

type PrependFilter struct {
	value string
}

func (a *PrependFilter) Prepend(input interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return a.value + typed
	default:
		return a.value + string(core.ToBytes(input))
	}
}
