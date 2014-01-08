package tags

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/filters"
	"testing"
)

func init() {
	core.RegisterFilter("minus", filters.MinusFactory)
}

func TestAssignForAStaticAndNoFilters(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" var = 'abc123'%}B")
	assertStringAssign(t, parser, "var", "abc123")
	spec.Expect(parser.Current()).ToEqual(byte('B'))
}

func TestAssignForAStaticWithFilters(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("sale  =  213  |minus: 4  %}o")
	assertIntAssign(t, parser, "sale", 209)
	spec.Expect(parser.Current()).ToEqual(byte('o'))
}

func TestAssignForAVariableWithFilters(t *testing.T) {
	parser := newParser("sale = price  |minus: 9  %}o")
	assertIntAssign(t, parser, "sale", 91)
}

func assertStringAssign(t *testing.T, parser *core.Parser, variableName, value string) {
	spec := gspec.New(t)
	tag, err := AssignFactory(parser, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("assign")
	m := make(map[string]interface{})
	tag.Execute(nil, m)
	spec.Expect(m[variableName].(string)).ToEqual(value)
}

func assertIntAssign(t *testing.T, parser *core.Parser, variableName string, value int) {
	spec := gspec.New(t)
	tag, err := AssignFactory(parser, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("assign")
	m := map[string]interface{}{
		"price": 100,
	}
	tag.Execute(nil, m)
	spec.Expect(m[variableName].(int)).ToEqual(value)
}
