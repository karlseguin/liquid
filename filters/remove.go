package filters

import (
	"github.com/karlseguin/liquid/core"
)

func RemoveFactory(parameters []core.Value) Filter {
	if len(parameters) != 1 {
		return Noop
	}
	return (&ReplaceFilter{parameters[0], EmptyValue, -1}).Replace
}
