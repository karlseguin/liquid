package filters

import (
	"bytes"
	"strings"
)

// Creates a upcase filter
func UpcaseFactory(parameters []string) Filter {
	return Upcase
}

// convert an input string to uppercase
func Upcase(input interface{}) interface{} {
	switch typed := input.(type) {
	case []byte:
		return bytes.ToUpper(typed)
	case string:
		return strings.ToUpper(typed)
	default:
		return input
	}
}
