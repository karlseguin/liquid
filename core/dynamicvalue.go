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
	for _, field := range v.Fields {
		if d = Resolve(d, field); d == nil {
			return nil
		}
	}
	return d
}

func (v *DynamicValue) Underlying() interface{} {
	return nil
}
