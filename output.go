package liquid

import (
	"fmt"
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/filters"
)

type Output struct {
	Value   core.Value
	Filters []filters.Filter
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

func newOutput(parser *core.Parser) (core.Code, error) {
	parser.ForwardBy(2) // skip the {{
	value, err := parser.ReadValue()
	if err != nil || value == nil {
		return nil, err
	}

	var filters []filters.Filter
	if parser.SkipSpaces() == '|' {
		parser.Forward()
		var err error
		if filters, err = buildFilters(parser); err != nil {
			return nil, err
		}
	}
	parser.ForwardBy(2) // skip the }}
	return &Output{value, filters}, nil
}

func buildFilters(p *core.Parser) ([]filters.Filter, error) {
	start := p.Position
	filters := make([]filters.Filter, 0, 5)
	for name := p.ReadName(); name != ""; name = p.ReadName() {
		factory, exists := Filters[name]
		if exists == false {
			return nil, p.Error(fmt.Sprintf("Unknown filter %q", name), start)
		}
		var parameters []core.Value
		if p.SkipSpaces() == ':' {
			p.Forward()
			var err error
			parameters, err = p.ReadParameters()
			if err != nil {
				return nil, err
			}
		}
		filters = append(filters, factory(parameters))
		if p.SkipSpaces() == '|' {
			p.Forward()
			continue
		}
		break
	}
	p.Commit()
	return filters, nil
}
