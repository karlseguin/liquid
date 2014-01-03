package core

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestParserToMarkupWhenTheresNoMarkup(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("hello world")
	pre, mt := parser.ToMarkup()
	spec.Expect(mt).ToEqual(0)
	spec.Expect(string(pre)).ToEqual("hello world")
	spec.Expect(parser.HasMore()).ToEqual(false)
}

func TestParserToMarkupWhenThereIsAnOutputMarkup(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("hello world {{ hello }}")
	pre, mt := parser.ToMarkup()
	spec.Expect(mt).ToEqual(OutputMarkup)
	spec.Expect(string(pre)).ToEqual("hello world ")
	spec.Expect(parser.HasMore()).ToEqual(true)
	spec.Expect(parser.Position).ToEqual(12)
}

func TestParserToMarkupWhenThereIsATagMarkup(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("hello world {% hello %}")
	pre, mt := parser.ToMarkup()
	spec.Expect(mt).ToEqual(TagMarkup)
	spec.Expect(string(pre)).ToEqual("hello world ")
	spec.Expect(parser.HasMore()).ToEqual(true)
	spec.Expect(parser.Position).ToEqual(12)
}

func TestParserSkipsSpacesWhenThereAreNoSpaces(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("hello")
	parser.SkipSpaces()
	spec.Expect(parser.Position).ToEqual(0)
}

func TestParserSkipsSpacesWhenThereAreSpaces(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("  hello")
	parser.SkipSpaces()
	spec.Expect(parser.Position).ToEqual(2)
}

func TestParserParsesAnEmptyValue(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("  ")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(value).ToBeNil()
	spec.Expect(parser.Position).ToEqual(2)
}

func TestParserParsesAnEmptyValue2(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("  }}")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(value).ToBeNil()
	spec.Expect(parser.Position).ToEqual(2)
}

func TestParserParsesAStaticValue(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello' ")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(string(value.Resolve(nil).([]byte))).ToEqual("hello")
	spec.Expect(parser.Position).ToEqual(8)
}

func TestParserParsesAnInteger(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 938 ")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(value.Resolve(nil).(int)).ToEqual(938)
	spec.Expect(parser.Position).ToEqual(4)
}

func TestParserParsesANegativeInteger(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" -331 ")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(value.Resolve(nil).(int)).ToEqual(-331)
	spec.Expect(parser.Position).ToEqual(5)
}

func TestParserParsesAFloat(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 9000.1 ")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(value.Resolve(nil).(float64)).ToEqual(9000.1)
	spec.Expect(parser.Position).ToEqual(7)
}

func TestParserParsesANegativeFloat(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" -331.89 ")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(value.Resolve(nil).(float64)).ToEqual(-331.89)
	spec.Expect(parser.Position).ToEqual(8)
}

func TestParserParsesAStaticValueWithEscapedQuote(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello \\'You\\' ' ")
	value, err := parser.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(string(value.Resolve(nil).([]byte))).ToEqual("hello 'You' ")
	spec.Expect(parser.Position).ToEqual(17)
}

func TestParserParsesAStaticWithMissingClosingQuote(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello ")
	_, err := parser.ReadValue()
	spec.Expect(err.Error()).ToEqual(`Invalid value, a single quote might be missing (" 'hello " - line 1)`)
}

func TestParserParsesASingleLevelDynamicValue(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" user ")
	v, err := parser.ReadValue()
	values := v.(*DynamicValue)
	spec.Expect(err).ToBeNil()
	spec.Expect(len(values.Fields)).ToEqual(1)
	spec.Expect(values.Fields[0]).ToEqual("user")
	spec.Expect(parser.Position).ToEqual(5)
}

func TestParserParsesAMultiLevelDynamicValue(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" user.NaMe.first}}")
	v, err := parser.ReadValue()
	values := v.(*DynamicValue)
	spec.Expect(err).ToBeNil()
	spec.Expect(len(values.Fields)).ToEqual(3)
	spec.Expect(values.Fields[0]).ToEqual("user")
	spec.Expect(values.Fields[1]).ToEqual("name")
	spec.Expect(values.Fields[2]).ToEqual("first")
	spec.Expect(parser.Position).ToEqual(16)
}

func TestParserReadsAnEmptyName1(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("  ")
	spec.Expect(parser.ReadName()).ToEqual("")
	spec.Expect(parser.Position).ToEqual(2)
}

func TestParserReadsAnEmptyName2(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("   }}")
	spec.Expect(parser.ReadName()).ToEqual("")
	spec.Expect(parser.Position).ToEqual(3)
}

func TestParserReadsAnEmptyName3(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("%}")
	spec.Expect(parser.ReadName()).ToEqual("")
	spec.Expect(parser.Position).ToEqual(0)
}

func TestParserReadsAnEmptyName4(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" |")
	spec.Expect(parser.ReadName()).ToEqual("")
	spec.Expect(parser.Position).ToEqual(1)
}

func TestParserReadsAName(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" spice }}")
	spec.Expect(parser.ReadName()).ToEqual("spice")
	spec.Expect(parser.Position).ToEqual(6)
}

func TestParserReadsEmptyParameters(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" }}")
	values, err := parser.ReadParameters()
	spec.Expect(err).ToBeNil()
	spec.Expect(len(values)).ToEqual(0)
	spec.Expect(parser.Position).ToEqual(1)
}

func TestParserReadsASingleParameter(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello'")
	values, err := parser.ReadParameters()
	spec.Expect(err).ToBeNil()
	spec.Expect(len(values)).ToEqual(1)
	spec.Expect(string(values[0].Resolve(nil).([]byte))).ToEqual("hello")
	spec.Expect(parser.Position).ToEqual(8)
}

func TestParserReadsMultipleParameters(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello' , 123 ")
	values, err := parser.ReadParameters()
	spec.Expect(err).ToBeNil()
	spec.Expect(len(values)).ToEqual(2)
	spec.Expect(string(values[0].Resolve(nil).([]byte))).ToEqual("hello")
	spec.Expect(values[1].Resolve(nil).(int)).ToEqual(123)
	spec.Expect(parser.Position).ToEqual(15)
}

func newParser(s string) *Parser {
	return NewParser([]byte(s))
}
