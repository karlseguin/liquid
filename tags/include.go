package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

// Creates an include tag
func IncludeFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	includeName, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	if includeName == nil {
		return nil, p.Error("Invalid include value")
	}

	valuesMap := map[string]core.Value{}

	next := p.SkipSpaces()
	if next == 'w' {
		name := p.ReadName()
		if name == "with" {
			value, _ := p.ReadValue()
			includeNameString := core.ToString(includeName.Resolve(nil))
			valuesMap[includeNameString] = value
		}
	}
	for {
		if next != ',' {
			break
		}
		p.Forward()

		// has one-line
		name := p.ReadName()

		next = p.SkipSpaces()
		if next == ':' {
			p.Forward()
			value, err := p.ReadValue()
			if err != nil {
				continue
			}
			valuesMap[name] = value
		}

		next = p.SkipSpaces()

	}

	p.SkipPastTag()

	return &Include{
		includeName: includeName,
		handler:     config.GetIncludeHandler(),
		valuesMap:   valuesMap,
	}, nil
}

type Include struct {
	includeName core.Value
	handler     core.IncludeHandler
	valuesMap   map[string]core.Value
}

func (i *Include) AddCode(code core.Code) {
	panic("Addcode should not have been called on a Include")
}

func (i *Include) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Include")
}

func (i *Include) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {

	// merge local values
	for key, val := range i.valuesMap {
		data[key] = val.Resolve(data)
	}

	if i.handler != nil {
		i.handler(core.ToString(i.includeName.Resolve(data)), writer, data)
	}

	return core.Normal
}

func (i *Include) Name() string {
	return "include"
}

func (i *Include) Type() core.TagType {
	return core.StandaloneTag
}
