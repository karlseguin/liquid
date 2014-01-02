package filters

import (
	"strconv"
)

var defaultMinus = &IntPlusFilter{-1, ""}

// Creates a minus filter
func MinusFactory(parameters []string) Filter {
	if len(parameters) == 0 || parameters[0] == "1" {
		return defaultMinus.Plus
	}
	if i, err := strconv.Atoi(parameters[0]); err == nil {
		return (&IntPlusFilter{-i, ""}).Plus
	}
	if f, err := strconv.ParseFloat(parameters[0], 64); err == nil {
		return (&FloatPlusFilter{-f, parameters[0]}).Plus
	}
	return defaultMinus.Plus
}
