package filters

import (
	"bytes"
	"github.com/karlseguin/liquid/core"
	"reflect"
)

var defaultJoin = (&JoinFilter{&core.StaticStringValue{" "}}).Join

// Creates a join filter
func JoinFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return defaultJoin
	}
	return (&JoinFilter{parameters[0]}).Join
}

type JoinFilter struct {
	glue core.Value
}

// join elements of the array with certain character between them
func (f *JoinFilter) Join(input interface{}, data map[string]interface{}) interface{} {
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
		array[i] = core.ToBytes(value.Index(i).Interface())
	}

	glue := core.ToBytes(f.glue.Resolve(data))
	return bytes.Join(array, glue)
}
