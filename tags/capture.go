package tags

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/core"
)

var endCapture = &End{"capture"}

func EndCaptureFactory(p *core.Parser) (core.Tag, error) {
	return endCapture, nil
}

// Creates an assign tag
func CaptureFactory(p *core.Parser) (core.Tag, error) {
	start := p.Position
	name := p.ReadName()
	if len(name) == 0 {
		return nil, p.Error("Invalid assignment, variable not found. ", start)
	}
	p.SkipPastTag()
	return &Capture{name, NewCommon()}, nil
}

type Capture struct {
	name string
	*Common
}

func (c *Capture) AddSibling(tag core.Tag) error {
	return errors.New(fmt.Sprintf("%q tag does not belong directly within a capture", tag.Name()))
}

func (c *Capture) Render(data map[string]interface{}) []byte {
	data[c.name] = c.Common.Render(data)
	return nil
}

func (c *Capture) Name() string {
	return "capture"
}

func (c *Capture) Type() core.TagType {
	return core.ContainerTag
}
