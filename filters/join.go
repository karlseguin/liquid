package filters

import (
	"bytes"
	"github.com/karlseguin/liquid/helpers"
	"reflect"
)

var defaultJoin = &JoinFilter{[]byte(" ")}

// Creates a join filter
func JoinFactory(parameters []string) Filter {
	if len(parameters) == 0 {
		return defaultJoin.Join
	}
	return (&JoinFilter{[]byte(parameters[0])}).Join
}

type JoinFilter struct {
	glue []byte
}

// join elements of the array with certain character between them
func (f *JoinFilter) Join(input interface{}) interface{} {
	value := reflect.ValueOf(input)
	kind := value.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return input
	}
	l := value.Len()
	if l == 0 {
		return input
	}

	array := make([][]byte, l)
	for i := 0; i < l; i++ {
		array[i] = helpers.ToBytes(value.Index(i).Interface())
	}

	return bytes.Join(array, f.glue)
}
