package filters

import (
	"strings"
)

// Creates a upcase filter
func UpcaseFactory(parameters []string) Filter {
	return Upcase
}

// convert an input string to uppercase
func Upcase(input interface{}) string {
  return strings.ToUpper(input.(string))
}
