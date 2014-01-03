package filters

import (
	"regexp"
)

// Used by other filters
type ReplacePattern struct {
	pattern *regexp.Regexp
	with    string
}

func (r *ReplacePattern) Replace(input interface{}, data map[string]interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return r.pattern.ReplaceAllLiteralString(typed, r.with)
	case []byte:
		return r.pattern.ReplaceAllLiteral(typed, []byte(r.with))
	default:
		return input
	}
}
