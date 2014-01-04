package filters

import (
	"bytes"
	"github.com/karlseguin/liquid/core"
	"strings"
)

func ReplaceFactory(parameters []core.Value) core.Filter {
	if len(parameters) != 2 {
		return Noop
	}
	return (&ReplaceFilter{parameters[0], parameters[1], -1}).Replace
}

type ReplaceFilter struct {
	target core.Value
	with   core.Value
	count  int
}

func (r *ReplaceFilter) Replace(input interface{}, data map[string]interface{}) interface{} {
	target := core.ToBytes(r.target.Resolve(data))
	with := core.ToBytes(r.with.Resolve(data))

	switch typed := input.(type) {
	case string:
		return strings.Replace(typed, string(target), string(with), r.count)
	case []byte:
		return bytes.Replace(typed, target, with, r.count)
	default:
		return input
	}
}
