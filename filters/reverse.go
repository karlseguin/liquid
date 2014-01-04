package filters

import (
	"github.com/karlseguin/liquid/core"
	"reflect"
)

// Creates a reverse filter
func ReverseFactory(parameters []core.Value) core.Filter {
	return Reverse
}

// reverses an array
func Reverse(input interface{}, data map[string]interface{}) interface{} {
	if s, ok := input.(string); ok {
		return reverseString(s)
	}

	value := reflect.ValueOf(input)
	kind := value.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return input
	}
	l := value.Len()
	if l < 2 {
		return input
	}

	mid := l / 2
	l -= 1
	for i := 0; i < mid; i++ {
		x := value.Index(l - i).Interface()
		value.Index(l - i).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(x))
	}
	return input
}

func reverseString(s string) []byte {
	b := []byte(s)
	l := len(b)
	if l < 2 {
		return b
	}
	mid := l / 2
	l -= 1
	for i := 0; i < mid; i++ {
		b[i], b[l-i] = b[l-i], b[i]
	}
	return b
}
