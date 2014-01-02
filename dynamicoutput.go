package liquid

import (
	"github.com/karlseguin/liquid/filters"
	"github.com/karlseguin/liquid/helpers"
	"strings"
)

type DynamicOutput struct {
	Fields  []string
	Filters []filters.Filter
}

func (o *DynamicOutput) Render(data interface{}) []byte {
	for _, field := range o.Fields {
		if data = helpers.Resolve(data, field); data == nil {
			return []byte("{{" + strings.Join(o.Fields, ".") + "}}")
		}
	}

	value := helpers.ResolveFinal(data)
	if o.Filters != nil {
		for _, filter := range o.Filters {
			value = filter(value)
		}
	}
	return helpers.ToBytes(value)
}

func createDynamicOutput(data, all []byte) (*DynamicOutput, int) {
	i := 0
	start := 0
	fields := make([]string, 0)
	for l := len(data); i < l; i++ {
		b := data[i]
		if b == ' ' {
			fields = append(fields, strings.ToLower(string(data[start:i])))
			break
		}
		if b == '.' {
			fields = append(fields, strings.ToLower(string(data[start:i])))
			start = i + 1
		}
	}
	return &DynamicOutput{
		Fields: helpers.TrimStrings(fields),
	}, i
}
