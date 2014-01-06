package filters

import (
	"github.com/karlseguin/liquid/core"
)

var (
	EmptyValue = &core.StaticStringValue{""}
)

func RemoveFirstFactory(parameters []core.Value) core.Filter {
	if len(parameters) != 1 {
		return Noop
	}
	return (&ReplaceFilter{parameters[0], EmptyValue, 1}).Replace
}
