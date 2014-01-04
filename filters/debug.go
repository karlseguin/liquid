package filters

import (
	"github.com/karlseguin/liquid/core"
	"strings"
)

func DebugFactory(parameter []core.Value) core.Filter {
	debug := &DebugFilter{parameter}
	return debug.Debug
}

type DebugFilter struct {
	parameters []core.Value
}

func (f *DebugFilter) Debug(input interface{}, data map[string]interface{}) interface{} {
	l := len(f.parameters)
	if l == 0 {
		return []byte("debug(" + input.(string) + ")")
	}
	values := make([]string, l)
	for i := 0; i < l; i++ {
		values[i] = string(core.ToBytes(f.parameters[i].Resolve(data)))
	}
	return []byte("debug(" + input.(string) + ", " + strings.Join(values, ", ") + ")")
}
