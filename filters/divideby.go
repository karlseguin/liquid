package filters

import (
	"strconv"
)

// Creates a divideby filter
func DivideByFactory(parameters []string) Filter {
	if len(parameters) == 0 {
		return Noop
	}
	if i, err := strconv.Atoi(parameters[0]); err == nil {
		return (&FloatTimesFilter{1 / float64(i)}).Times
	}
	if f, err := strconv.ParseFloat(parameters[0], 64); err == nil {
		return (&FloatTimesFilter{1 / f}).Times
	}
	return Noop
}
