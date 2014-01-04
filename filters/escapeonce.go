package filters

import (
	"github.com/karlseguin/liquid/core"
	"html"
)

// creates an escapeonce filter
func EscapeOnceFactory(parameters []core.Value) core.Filter {
	return EscapeOnce
}

// Escapes html that hasn't already been escaped
func EscapeOnce(input interface{}, data map[string]interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return html.EscapeString(html.UnescapeString(typed))
	case []byte:
		return html.EscapeString(html.UnescapeString(string(typed)))
	default:
		return input
	}
}
