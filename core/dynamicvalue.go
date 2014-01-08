package core

import (
	"reflect"
	"strings"
)

type DynamicValue struct {
	Fields []string
	index  Value
}

func NewDynamicValue(fields []string) *DynamicValue {
	value := &DynamicValue{fields, nil}
	last := len(fields) - 1
	if index, field, ok := unindexDynamicField(fields[last]); ok {
		fields[last] = field
		value.index = index
	}
	return value
}

func (v *DynamicValue) ResolveWithNil(data map[string]interface{}) interface{} {
	d := v.resolve(data)
	if d == nil {
		return nil
	}
	return ResolveFinal(d)
}

func (v *DynamicValue) Resolve(data map[string]interface{}) interface{} {
	d := v.resolve(data)
	if d == nil {
		return []byte("{{" + strings.Join(v.Fields, ".") + "}}")
	}
	return ResolveFinal(d)
}

func (v *DynamicValue) resolve(data map[string]interface{}) interface{} {
	var d interface{} = data
	var p interface{}

	for i, l := 0, len(v.Fields); i < l; i++ {
		field := v.Fields[i]
		if d = Resolve(d, field); d == nil {
			if field == "size" && i == l-1 && p != nil {
				if n, ok := ToLength(p); ok {
					return n
				}
			}
			return nil
		}
		p = d
	}
	if v.index != nil {
		return indexedValue(d, v.index.Resolve(data))
	}
	return d
}

func (v *DynamicValue) Underlying() interface{} {
	return nil
}

func unindexDynamicField(field string) (Value, string, bool) {
	end := len(field) - 1
	if field[end] != ']' {
		return nil, field, false
	}

	start := end
	for ; start >= 0; start-- {
		if field[start] == '[' {
			break
		}
	}
	if start == 0 {
		return nil, field, false
	}
	value, err := NewParser([]byte(field[start+1:end] + " ")).ReadValue()
	if err != nil {
		return nil, field, false
	}
	return value, field[0:start], true
}

func indexedValue(container interface{}, index interface{}) interface{} {
	value := reflect.ValueOf(container)
	kind := value.Kind()
	if kind == reflect.Array || kind == reflect.Slice || kind == reflect.String {
		if n, ok := ToInt(index); ok && n <= value.Len() && n > 0 {
			return value.Index(n - 1).Interface()
		}
	} else if kind == reflect.Map {
		indexValue := reflect.ValueOf(index)
		return value.MapIndex(indexValue).Interface()
	}
	return nil
}
