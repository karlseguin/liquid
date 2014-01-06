package tags

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/core"
	"io"
)

var endCapture = &End{"capture"}

func EndCaptureFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endCapture, nil
}

// Creates an assign tag
func CaptureFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	start := p.Position
	name := p.ReadName()
	if len(name) == 0 {
		return nil, p.Error("Invalid assignment, variable not found. ", start)
	}
	p.SkipPastTag()
	return &Capture{name, config, NewCommon()}, nil
}

type Capture struct {
	name   string
	config *core.Configuration
	*Common
}

func (c *Capture) AddSibling(tag core.Tag) error {
	return errors.New(fmt.Sprintf("%q tag does not belong directly within a capture", tag.Name()))
}

func (c *Capture) Render(w io.Writer, data map[string]interface{}) {
	writer := c.config.GetWriter()
	defer writer.Close()
	c.Common.Render(writer, data)
	data[c.name] = writer.Bytes()
}

func (c *Capture) Name() string {
	return "capture"
}

func (c *Capture) Type() core.TagType {
	return core.ContainerTag
}
