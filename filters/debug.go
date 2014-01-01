package filters

import (
	"strings"
)

func DebugFactory(parameter []string) Filter {
	debug := &DebugFilter{parameter}
	return debug.Debug
}

type DebugFilter struct {
	parameters []string
}

func (f *DebugFilter) Debug(input interface{}) interface{} {
	if len(f.parameters) == 0 {
		return []byte("debug(" + input.(string) + ")")
	}
	return []byte("debug(" + input.(string) + ", " + strings.Join(f.parameters, ", ") + ")")
}
