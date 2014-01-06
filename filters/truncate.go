package filters

import (
	"github.com/karlseguin/liquid/core"
)

var (
	defaultTruncateLimit  = &core.StaticIntValue{50}
	defaultTruncateAppend = &core.StaticStringValue{"..."}
	defaultTruncate       = (&TruncateFilter{defaultTruncateLimit, defaultTruncateAppend}).Truncate
)

// Creates an truncate filter
func TruncateFactory(parameters []core.Value) core.Filter {
	switch len(parameters) {
	case 0:
		return defaultTruncate
	case 1:
		return (&TruncateFilter{parameters[0], defaultTruncateAppend}).Truncate
	default:
		return (&TruncateFilter{parameters[0], parameters[1]}).Truncate
	}
}

type TruncateFilter struct {
	limit  core.Value
	append core.Value
}

func (t *TruncateFilter) Truncate(input interface{}, data map[string]interface{}) interface{} {
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

	append := core.ToString(t.append.Resolve(data))
	length -= len(append)
	if length >= len(value) {
		return input
	}
	if length < 0 {
		return append
	}
	return value[:length] + append
}
