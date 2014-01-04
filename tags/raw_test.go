package tags

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TEstRawFactory(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} this {{}} {%} is raw {%endraw%}Z")
	tag, err := RawFactory(parser)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("raw")
	spec.Expect(parser.Current()).ToEqual(byte('Z'))
}

func TestRawFactoryHandlesUnclosedRaw(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} this is raw {%enccsad%}X")
	tag, err := RawFactory(parser)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("raw")
	spec.Expect(parser.HasMore()).ToEqual(false)
}

func TestRawTagRenders(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} this {{}} {%} is raw {%endraw%}Z")
	tag, _ := RawFactory(parser)
	spec.Expect(string(tag.Render(nil))).ToEqual(" this {{}} {%} is raw ")
}
