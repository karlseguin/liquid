package filters

import (
	"reflect"
)

// Creates a first filter
func FirstFactory(parameters []string) Filter {
	return First
}

// get the first element of the passed in array
func First(input interface{}) interface{} {
	value := reflect.ValueOf(input)
	kind := value.Kind()
	if (kind != reflect.Array && kind != reflect.Slice) || value.Len() == 0 {
		return input
	}
	return value.Index(0).Interface()
}
