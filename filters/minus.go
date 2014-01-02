package filters

import (
	"strconv"
)

// Creates a minus filter
func MinusFactory(parameters []string) Filter {
	if len(parameters) == 0 {
		return (&IntPlusFilter{-1, ""}).Plus
	}
	if i, err := strconv.Atoi(parameters[0]); err == nil {
		return (&IntPlusFilter{-i, ""}).Plus
	}
	return (&IntPlusFilter{0, ""}).Plus
}
