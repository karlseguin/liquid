package filters

import (
	"bytes"
	"strings"
)

// Creates a downcase filter
func DowncaseFactory(parameters []string) Filter {
	return Downcase
}

// convert an input string to lowercase
func Downcase(input interface{}) interface{} {
	switch typed := input.(type) {
	case []byte:
		return bytes.ToUpper(typed)
	case string:
		return strings.ToUpper(typed)
	default:
		return input
	}
}
