package filters

import (
	"github.com/karlseguin/liquid/core"
)

func ReplaceFirstFactory(parameters []core.Value) core.Filter {
	if len(parameters) != 2 {
		return Noop
	}
	return (&ReplaceFilter{parameters[0], parameters[1], 1}).Replace
}
