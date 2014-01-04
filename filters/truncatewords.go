package filters

import (
	"github.com/karlseguin/liquid/core"
)

var (
	defaultTruncateWordsLimit = &core.StaticIntValue{15}
	defaultTruncateWords      = (&TruncateWordsFilter{defaultTruncateWordsLimit, defaultTruncateAppend}).Truncate
)

// Creates an truncatewords filter
func TruncateWordsFactory(parameters []core.Value) core.Filter {
	switch len(parameters) {
	case 0:
		return defaultTruncateWords
	case 1:
		return (&TruncateWordsFilter{parameters[0], defaultTruncateAppend}).Truncate
	default:
		return (&TruncateWordsFilter{parameters[0], parameters[1]}).Truncate
	}
}

type TruncateWordsFilter struct {
	limit  core.Value
	append core.Value
}

func (t *TruncateWordsFilter) Truncate(input interface{}, data map[string]interface{}) interface{} {
	length, ok := core.ToInt(t.limit.Resolve(data))
	if ok == false {
		return input
	}

	var value string
	switch typed := input.(type) {
	case string:
		value = typed
	default:
		value = core.ToString(typed)
	}

	i := 2
	found := 0
	l := len(value)
	for ; i < l; i++ {
		if value[i] == ' ' && value[i-1] != ' ' && value[i-2] != ' ' {
			found++
			if found == length {
				break
			}
		}
	}

	if i == l {
		return input
	}
	return value[:i] + core.ToString(t.append.Resolve(data))
}
