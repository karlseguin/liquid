package filters

import (
	"github.com/karlseguin/liquid/helpers"
)

// Creates an append filter
func AppendFactory(parameters []string) Filter {
	if len(parameters) == 0 || parameters[0] == "" {
		return Noop
	}
	return (&AppendFilter{parameters[0]}).Append
}

type AppendFilter struct {
	value string
}

func (a *AppendFilter) Append(input interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return typed + a.value
	default:
		return string(helpers.ToBytes(input)) + a.value
	}
}
