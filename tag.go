package liquid

import (
	"fmt"
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/tags"
)

type TagFactory func(*core.Parser) (core.Tag, error)

var Tags = map[string]TagFactory{
	"comment":    tags.CommentFactory,
	"endcomment": tags.EndCommentFactory,
	"raw":        tags.RawFactory,
	"endraw":     tags.EndRawFactory,
	"assign":     tags.AssignFactory,
	"capture":    tags.CaptureFactory,
	"endcapture": tags.EndCaptureFactory,
	"if":         tags.IfFactory,
	"elseif":     tags.ElseIfFactory,
	"else":       tags.ElseFactory,
	"endif":      tags.EndIfFactory,
	"unless":     tags.UnlessFactory,
	"endunless":  tags.EndUnlessFactory,
	"case":       tags.CaseFactory,
	"when":       tags.WhenFactory,
	"endcase":    tags.EndCaseFactory,
}

func newTag(p *core.Parser) (core.Tag, error) {
	start := p.Position
	p.ForwardBy(2) // skip the {%
	name := p.ReadName()
	factory, ok := Tags[name]
	if ok == false {
		return nil, p.Error(fmt.Sprintf("unknown tag %q", name), start)
	}
	return factory(p)
}
