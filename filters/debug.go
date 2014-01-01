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

func (f *DebugFilter) Debug(input interface{}) string {
	if len(f.parameters) == 0 {
		return "debug(" + input.(string) + ")"
	}
	return "debug(" + input.(string) + ", " + strings.Join(f.parameters, ", ") + ")"
}
