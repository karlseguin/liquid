package tags

import (
	"github.com/karlseguin/liquid/core"
)

var theOneComment = new(Comment)

func CommentFactory(data, all []byte) (core.Token, error) {
	return theOneComment, nil
}

type Comment struct {
}

func (c *Comment) Name() string {
	return "comment"
}

func (c *Comment) IsEnd() bool {
	return false
}

func (c *Comment) Render(data map[string]interface{}) []byte {
	return []byte{}
}

func (c *Comment) AddToken(token core.Token) {
}

func (c *Comment) AddTag(tag core.Tag) (bool, bool) {
	if tag.Name() == "endcomment" {
		return true, false
	}
	return false, false
}
