package tags

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestCommentFactoryForNormalComment(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} hack {%endcomment%}Z")
	tag, err := CommentFactory(parser, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("comment")
	spec.Expect(parser.Current()).ToEqual(byte('Z'))
}

func TestCommentFactoryForNestedComment(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} ha {%comment%} {%if%} ck {%endcomment%} {%  endcomment  %}XZ ")
	tag, err := CommentFactory(parser, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("comment")
	spec.Expect(parser.Current()).ToEqual(byte('X'))
}

func TestCommentFactoryHandlesUnclosedComment(t *testing.T) {
	spec := gspec.New(t)
	parser := newParser(" %} ouch ")
	tag, err := CommentFactory(parser, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(tag.Name()).ToEqual("comment")
	spec.Expect(parser.HasMore()).ToEqual(false)
}

func newParser(s string) *core.Parser {
	return core.NewParser([]byte(s))
}
