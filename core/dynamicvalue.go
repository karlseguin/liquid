package core

import (
	"strings"
)

type DynamicValue struct {
	Fields []string
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

	for i, l := 0, len(v.Fields); i < l; i ++ {
		field := v.Fields[i]
		if d = Resolve(d, field); d == nil {
			if field == "size" && i == l - 1 && p != nil {
				if n, ok := ToLength(p); ok {
					return n
				}
			}
			return nil
		}
		p = d
	}
	return d
}

func (v *DynamicValue) Underlying() interface{} {
	return nil
}
