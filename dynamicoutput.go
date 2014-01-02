package liquid

import (
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/filters"
	"strings"
)

type DynamicOutput struct {
	Fields  []string
	Filters []filters.Filter
}

func (o *DynamicOutput) Render(data map[string]interface{}) []byte {
	var d interface{} = data
	for _, field := range o.Fields {
		if d = core.Resolve(d, field); data == nil {
			return []byte("{{" + strings.Join(o.Fields, ".") + "}}")
		}
	}

	value := core.ResolveFinal(d)
	if o.Filters != nil {
		for _, filter := range o.Filters {
			value = filter(value)
		}
	}
	return core.ToBytes(value)
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
		Fields: core.TrimStrings(fields),
	}, i
}
