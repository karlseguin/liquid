package liquid

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

type Output struct {
	Value   core.Value
	Filters []core.Filter
}

func (o *Output) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	value := o.Value.Resolve(data)
	if o.Filters != nil {
		for _, filter := range o.Filters {
			value = filter(value, data)
		}
	}
	writer.Write(core.ToBytes(value))
	return core.Normal
}

func newOutput(p *core.Parser) (core.Code, error) {
	p.ForwardBy(2) // skip the {{
	value, err := p.ReadValue()
	if err != nil || value == nil {
		return nil, err
	}

	filters, err := p.ReadFilters()
	if err != nil {
		return nil, err
	}
	p.SkipPastOutput()
	return &Output{value, filters}, nil
}
