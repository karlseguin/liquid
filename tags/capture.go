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
	name := p.ReadName()
	if len(name) == 0 {
		return nil, p.Error("Invalid assignment, variable not found. ")
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

func (c *Capture) Execute(w io.Writer, data map[string]interface{}) core.ExecuteState {
	writer := core.BytePool.Checkout()
	defer writer.Close()
	c.Common.Execute(writer, data)
	data[c.name] = writer.Bytes()
	return core.Normal
}

func (c *Capture) Name() string {
	return "capture"
}

func (c *Capture) Type() core.TagType {
	return core.ContainerTag
}
