package liquid

import (
	"github.com/karlseguin/liquid/core"
)

type Output struct {
	Value   core.Value
	Filters []core.Filter
}

func (o *Output) Render(data map[string]interface{}) []byte {
	value := o.Value.Resolve(data)
	if o.Filters != nil {
		for _, filter := range o.Filters {
			value = filter(value, data)
		}
	}
	return core.ToBytes(value)
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
