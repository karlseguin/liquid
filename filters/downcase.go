package filters

import (
	"bytes"
	"github.com/karlseguin/liquid/core"
	"strings"
)

// Creates a downcase filter
func DowncaseFactory(parameters []core.Value) core.Filter {
	return Downcase
}

// convert an input string to lowercase
func Downcase(input interface{}, data map[string]interface{}) interface{} {
	switch typed := input.(type) {
	case []byte:
		return bytes.ToLower(typed)
	case string:
		return strings.ToLower(typed)
	default:
		return input
	}
}
