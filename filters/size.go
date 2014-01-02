package filters

import (
	"reflect"
)

// Creates a size filter
func SizeFactory(parameters []string) Filter {
	return Size
}

// Gets the size of a string or array
func Size(input interface{}) interface{} {
	if s, ok := input.(string); ok {
		return len(s)
	}
	value := reflect.ValueOf(input)
	kind := value.Kind()
	if kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map {
		return value.Len()
	}
	return input
}
