package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

// Creates an include tag
func IncludeFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	value, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, p.Error("Invalid include value")
	}
	p.SkipPastTag()
	return &Include{value, config.GetIncludeHandler()}, nil
}

type Include struct {
	value   core.Value
	handler core.IncludeHandler
}

func (i *Include) AddCode(code core.Code) {
	panic("Addcode should not have been called on a Include")
}

func (i *Include) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Include")
}

func (i *Include) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	if i.handler != nil {
		i.handler(core.ToString(i.value.Resolve(data)), writer, data)
	}
	return core.Normal
}

func (i *Include) Name() string {
	return "include"
}

func (i *Include) Type() core.TagType {
	return core.StandaloneTag
}
