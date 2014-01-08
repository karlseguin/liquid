package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

var continueTag = new(Continue)

// Creates a continue tag
func ContinueFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	p.SkipPastTag()
	return continueTag, nil
}

type Continue struct{}

func (c *Continue) AddCode(code core.Code) {
	panic("AddCode should not have been called on a Continue")
}

func (c *Continue) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Continue")
}

func (c *Continue) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	return core.Continue
}

func (c *Continue) Name() string {
	return "continue"
}

func (c *Continue) Type() core.TagType {
	return core.StandaloneTag
}
