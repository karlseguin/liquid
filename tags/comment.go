package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

var comment = new(Comment)
var endComment = &End{"comment"}

// Special handling to just quickly skip over it all
func CommentFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	openTags := 1
	for {
		_, markupType := p.ToMarkup(false)
		if markupType == core.TagMarkup {
			p.ForwardBy(2) // skip {%
			if name := p.ReadName(); name == "comment" {
				openTags++
			} else if name == "endcomment" {
				openTags--
				if openTags == 0 {
					p.SkipPastTag()
					break
				}
			}
		} else if markupType == core.OutputMarkup {

			p.SkipPastTag()
		} else {
			break
		}
	}
	return comment, nil
}

func EndCommentFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endComment, nil
}

// Comment tag is a special tag in that, while it looks like a container tag,
// we treat it as an end tag and just move the parser all the way past the
// end tag. A
type Comment struct {
}

func (c *Comment) AddCode(code core.Code) {
	panic("AddCode should not have been called on a comment")
}

func (c *Comment) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a comment")
}

func (c *Comment) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	panic("Render should not have been called on a comment")
}

func (c *Comment) Name() string {
	return "comment"
}

func (c *Comment) Type() core.TagType {
	return core.NoopTag
}
