package filters

import (
	"bytes"
	"strings"
)

func ReplaceFactory(parameters []string) Filter {
	if len(parameters) != 2 {
		return Noop
	}
	return (&ReplaceFilter{parameters[0], parameters[1], -1}).Replace
}

type ReplaceFilter struct {
	target string
	with   string
	count  int
}

func (r *ReplaceFilter) Replace(input interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return strings.Replace(typed, r.target, r.with, r.count)
	case []byte:
		return bytes.Replace(typed, []byte(r.target), []byte(r.with), r.count)
	default:
		return input
	}
}
