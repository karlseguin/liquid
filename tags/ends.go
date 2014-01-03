package tags

import (
	"github.com/karlseguin/liquid/core"
)

var (
	EndComment = &EndTag{"endcomment"}
	EndRaw     = &EndTag{"endraw"}
)

type EndTag struct {
	name string
}

func (e *EndTag) Name() string {
	return e.name
}

func (e *EndTag) IsEnd() bool {
	return true
}

func (e *EndTag) Render(data map[string]interface{}) []byte {
	panic("Render should not be called on an end tag")
}

func (e *EndTag) AddToken(token core.Token) {
	panic("AddToken should not be called on an end tag")
}

func (e *EndTag) AddTag(tag core.Tag) (bool, bool) {
	panic("AddTag should not be called on an end tag")
}
