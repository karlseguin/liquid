package liquid

import (
	"github.com/karlseguin/liquid/filters"
)

type OutputStatic struct {
	Value   []byte
	Filters []filters.Filter
}

func (o *OutputStatic) Render(data interface{}) []byte {
	if o.Filters == nil {
		return o.Value
	}
	var value interface{} = o.Value
	for _, filter := range o.Filters {
		value = filter(value)
	}
	return value.([]byte)
}
