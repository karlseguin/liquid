package filters

import (
	"github.com/karlseguin/liquid/core"
)

func ReplaceFirstFactory(parameters []core.Value) Filter {
	if len(parameters) != 2 {
		return Noop
	}
	return (&ReplaceFilter{parameters[0], parameters[1], 1}).Replace
}
