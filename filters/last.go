package filters

import (
	"reflect"
)

// Creates a last filter
func LastFactory(parameters []string) Filter {
	return Last
}

// get the last element of the passed in array
func Last(input interface{}) interface{} {
	value := reflect.ValueOf(input)
	kind := value.Kind()

	if (kind != reflect.Array && kind != reflect.Slice) {
		return input
	}
	len := value.Len()
	if len == 0 {
		return input
	}
	return value.Index(len - 1).Interface()
}
