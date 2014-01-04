package filters

import (
	"github.com/karlseguin/liquid/core"
	"reflect"
)

var defaultDefaultFilter = (&DefaultFilter{EmptyValue}).Default

// Creates a default filter
func DefaultFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return defaultDefaultFilter
	}
	return (&DefaultFilter{parameters[0]}).Default
}

type DefaultFilter struct {
	value core.Value
}

func (d *DefaultFilter) Default(input interface{}, data map[string]interface{}) interface{} {
	dflt := d.value.Resolve(data)

	if input == nil {
		return dflt
	}

	value := reflect.ValueOf(input)
	kind := value.Kind()
	if kind != reflect.Array && kind != reflect.Slice && kind != reflect.Map && kind != reflect.String {
		return input
	}
	if value.Len() == 0 {
		return dflt
	}
	return input
}
