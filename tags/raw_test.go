package tags

import (
	"bytes"
	"github.com/karlseguin/gspec"
	"testing"
)

func TEstRawFactory(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} this {{}} {%} is raw {%endraw%}Z")
	tag, err := RawFactory(parser, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("raw")
	spec.Expect(parser.Current()).ToEqual(byte('Z'))
}

func TestRawFactoryHandlesUnclosedRaw(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} this is raw {%enccsad%}X")
	tag, err := RawFactory(parser, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("raw")
	spec.Expect(parser.HasMore()).ToEqual(false)
}

func TestRawTagExecutes(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} this {{}} {%} is raw {%endraw%}Z")
	tag, _ := RawFactory(parser, nil)

	writer := new(bytes.Buffer)
	tag.Execute(writer, nil)
	spec.Expect(writer.String()).ToEqual(" this {{}} {%} is raw ")
}
