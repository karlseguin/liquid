package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

var breakTag = new(Break)

// Creates a break tag
func BreakFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	p.SkipPastTag()
	return breakTag, nil
}

type Break struct{}

func (b *Break) AddCode(code core.Code) {
	panic("Addcode should not have been called on a Break")
}

func (b *Break) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Break")
}

func (b *Break) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	return core.Break
}

func (b *Break) Name() string {
	return "break"
}

func (b *Break) Type() core.TagType {
	return core.StandaloneTag
}
