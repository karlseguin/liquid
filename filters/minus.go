package filters

import (
	"strconv"
)

var defaultMinus = &IntPlusFilter{-1, ""}

// Creates a minus filter
func MinusFactory(parameters []string) Filter {
	if len(parameters) == 0 {
		return defaultMinus.Plus
	}
	if i, err := strconv.Atoi(parameters[0]); err == nil {
		return (&IntPlusFilter{-i, ""}).Plus
	}
	return defaultMinus.Plus
}
