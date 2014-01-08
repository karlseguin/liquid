package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

var endRaw = &End{"raw"}

// Special handling to just quickly skip over it all
func RawFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	p.SkipPastTag()
	start := p.Position
	end := start
	for {
		_, markupType := p.ToMarkup(false)
		if markupType == core.TagMarkup {
			//tentative end is before the start of the endraw tag
			end = p.Position
			p.ForwardBy(2) // skip {%
			if name := p.ReadName(); name == "endraw" {
				p.SkipPastTag()
				break
			}
		} else if markupType == core.OutputMarkup {
			p.ForwardBy(2) // skip it
		} else {
			break
		}
	}
	return &Raw{p.Data[start:end]}, nil
}

func EndRawFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endRaw, nil
}

// Raw tag is a special tag. Like comment it'll forward the parser to past its
// own end.
type Raw struct {
	value []byte
}

func (r *Raw) AddCode(code core.Code) {
	panic("AddCode should not have been called on a Raw")
}

func (r *Raw) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Raw")
}

func (r *Raw) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	writer.Write(r.value)
	return core.Normal
}

func (r *Raw) Name() string {
	return "raw"
}

func (r *Raw) Type() core.TagType {
	return core.StandaloneTag
}
