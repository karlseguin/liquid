package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

// a generic end tag
type End struct {
	name string
}

func (e *End) AddCode(code core.Code) {
	panic("AddCode should not have been called on an end tag")
}

func (e *End) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on an end tag")
}

func (e *End) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	panic("Render should not have been called on an end tag")
}

func (e *End) Name() string {
	return e.name
}

func (e *End) Type() core.TagType {
	return core.EndTag
}
