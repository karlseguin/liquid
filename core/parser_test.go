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
	spec.Expect(string(pre)).ToEqual("hello world")
	spec.Expect(parser.HasMore()).ToEqual(true)
	spec.Expect(parser.Peek()).ToEqual(byte(' '))
}

func TestParserToMarkupWhenThereIsATagMarkup(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("hello world {% hello %}")
	pre, mt := parser.ToMarkup()
	spec.Expect(mt).ToEqual(TagMarkup)
	spec.Expect(string(pre)).ToEqual("hello world")
	spec.Expect(parser.HasMore()).ToEqual(true)
	spec.Expect(parser.Peek()).ToEqual(byte(' '))
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
	v, static, err := parser.ReadValue()
	value := string(v.([]byte))
	spec.Expect(err).ToBeNil()
	spec.Expect(static).ToEqual(true)
	spec.Expect(value).ToEqual("")
	spec.Expect(parser.Position).ToEqual(2)
}

func TestParserParsesAnEmptyValue2(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser("  }}")
	v, static, err := parser.ReadValue()
	value := string(v.([]byte))
	spec.Expect(err).ToBeNil()
	spec.Expect(static).ToEqual(true)
	spec.Expect(value).ToEqual("")
	spec.Expect(parser.Position).ToEqual(2)
}

func TestParserParsesAStaticValue(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello' ")
	v, static, err := parser.ReadValue()
	value := string(v.([]byte))
	spec.Expect(err).ToBeNil()
	spec.Expect(static).ToEqual(true)
	spec.Expect(value).ToEqual("hello")
	spec.Expect(parser.Position).ToEqual(8)
}

func TestParserParsesAStaticValueWithEscapedQuote(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello \\'You\\' ' ")
	v, static, err := parser.ReadValue()
	value := string(v.([]byte))
	spec.Expect(err).ToBeNil()
	spec.Expect(static).ToEqual(true)
	spec.Expect(value).ToEqual("hello 'You' ")
	spec.Expect(parser.Position).ToEqual(17)
}

func TestParserParsesAStaticWithMissingClosingQuote(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" 'hello ")
	_, _, err := parser.ReadValue()
	spec.Expect(err.Error()).ToEqual(`Invalid output value, a single quote might be missing (" 'hello " - line 1)`)
}

func TestParserParsesASingleLevelDynamicValue(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" user ")
	v, static, err := parser.ReadValue()
	values := v.([]string)
	spec.Expect(err).ToBeNil()
	spec.Expect(static).ToEqual(false)
	spec.Expect(len(values)).ToEqual(1)
	spec.Expect(values[0]).ToEqual("user")
	spec.Expect(parser.Position).ToEqual(5)
}

func TestParserParsesAMultiLevelDynamicValue(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" user.NaMe.first}}")
	v, static, err := parser.ReadValue()
	values := v.([]string)
	spec.Expect(err).ToBeNil()
	spec.Expect(static).ToEqual(false)
	spec.Expect(len(values)).ToEqual(3)
	spec.Expect(values[0]).ToEqual("user")
	spec.Expect(values[1]).ToEqual("name")
	spec.Expect(values[2]).ToEqual("first")
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

func newParser(s string) *Parser {
	return NewParser([]byte(s))
}
