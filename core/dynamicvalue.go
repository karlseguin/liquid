package core

import (
	"fmt"
	"strings"
)

type DynamicValue struct {
	Fields []string
}

func (v *DynamicValue) Resolve(data map[string]interface{}) interface{} {
	var d interface{} = data
	for _, field := range v.Fields {
		if d = Resolve(d, field); d == nil {
			return []byte("{{" + strings.Join(v.Fields, ".") + "}}")
		}
	}

	return ResolveFinal(d)
}

func (v *DynamicValue) Underlying() interface{} {
	return nil
}
