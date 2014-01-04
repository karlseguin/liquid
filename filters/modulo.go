package filters

import (
	"github.com/karlseguin/liquid/core"
	"strconv"
)

// Creates a plus filter
func ModuloFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return Noop
	}
	switch typed := parameters[0].(type) {
	case *core.StaticIntValue:
		return (&IntModuloFilter{typed.Value}).Mod
	case *core.DynamicValue:
		return (&DynamicModuloFilter{typed}).Mod
	default:
		if n, ok := core.ToInt(parameters[0].Underlying()); ok {
			return (&IntModuloFilter{int(n)}).Mod
		}
		return Noop
	}
}

type DynamicModuloFilter struct {
	value core.Value
}

func (m *DynamicModuloFilter) Mod(input interface{}, data map[string]interface{}) interface{} {
	if n, ok := core.ToInt(m.value.Resolve(data)); ok {
		return modInt(n, input)
	}
	return input
}

// Plus filter for integer parameter
type IntModuloFilter struct {
	mod int
}

func (m *IntModuloFilter) Mod(input interface{}, data map[string]interface{}) interface{} {
	return modInt(m.mod, input)
}

func modInt(mod int, input interface{}) interface{} {
	switch typed := input.(type) {
	case int:
		return typed % mod
	case uint:
		return typed % uint(mod)
	case int64:
		return typed % int64(mod)
	case uint64:
		return typed % uint64(mod)
	case string:
		return stringModInt(typed, mod)
	case []byte:
		return stringModInt(string(typed), mod)
	default:
		return input
	}
}

func stringModInt(s string, mod int) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return i % mod
	}
	return s
}
