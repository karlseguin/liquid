package filters

import (
	"strings"
)

// Creates a downcase filter
func DowncaseFactory(parameters []string) Filter {
	return Downcase
}

// convert an input string to lowercase
func Downcase(input interface{}) string {
  return strings.ToLower(input.(string))
}
