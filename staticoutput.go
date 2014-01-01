package liquid

import (
	"github.com/karlseguin/liquid/filters"
	"fmt"
	"errors"
)

type StaticOutput struct {
	Value   []byte
	Filters []filters.Filter
}

func (o *StaticOutput) Render(data interface{}) []byte {
	if o.Filters == nil {
		return o.Value
	}
	var value interface{} = o.Value
	for _, filter := range o.Filters {
		value = filter(value)
	}
	return value.([]byte)
}


func createStaticOutput(data, all []byte) (*StaticOutput, int, error) {
	escaped := 0
	escaping := false
	for index, b := range data {
		if b == '\'' {
			if escaping {
				escaped++
				escaping = false
			} else {
				var value []byte
				if escaped > 0 {
					return &StaticOutput{Value: unescape(data[0:index], escaped)}, index, nil
				}
				value = make([]byte, index)
				copy(value, data[:index])
				return &StaticOutput{Value: value}, index, nil
			}
		} else if b == '\\' && escaping == false {
			escaping = true
		} else {
			escaping = false
		}
	}
	return nil, 0, errors.New(fmt.Sprintf("Output had an unclosed single quote in %q", all))
}

func unescape(data []byte, escaped int) []byte {
	value := make([]byte, len(data)-escaped)
	i := 0
	found := 0
	position := 0
	for l := len(data) - 1; i < l; i++ {
		b := data[i]
		if b == '\\' && data[i+1] == '\'' {
			value[position] = '\''
			found++
			i++
			if found == escaped {
				break
			}
		} else {
			value[position] = b
		}
		position++
	}
	copy(value[position:], data[i:])
	return value
}
