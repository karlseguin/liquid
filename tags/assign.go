package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

// Creates an assign tag
func AssignFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	start := p.Position
	name := p.ReadName()
	if len(name) == 0 {
		return nil, p.Error("Invalid variable name in assign tag", start)
	}
	if p.SkipUntil('=') != '=' {
		return nil, p.Error("Invalid assign, missing '=' ", start)
	}
	p.Forward()
	value, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	filters, err := p.ReadFilters()
	if err != nil {
		return nil, err
	}
	p.SkipPastTag()
	return &Assign{name, value, filters}, nil
}

type Assign struct {
	name    string
	value   core.Value
	filters []core.Filter
}

func (a *Assign) AddCode(code core.Code) {
	panic("Addcode should not have been called on a Assign")
}

func (a *Assign) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Assign")
}

func (a *Assign) Render(writer io.Writer, data map[string]interface{}) {
	value := a.value.Resolve(data)
	if a.filters != nil {
		for _, filter := range a.filters {
			value = filter(value, data)
		}
	}
	data[a.name] = value
}

func (a *Assign) Name() string {
	return "assign"
}

func (a *Assign) Type() core.TagType {
	return core.StandaloneTag
}
