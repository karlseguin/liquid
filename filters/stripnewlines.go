package filters

import (
	"regexp"
)

var stripNewLinesPattern = regexp.MustCompile("(\n|\r)")
var emptyBytes = []byte("")

func StripNewLinesFactory(parameters []string) Filter {
	return StripNewLines
}

func StripNewLines(input interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return stripNewLinesPattern.ReplaceAllString(typed, "")
	case []byte:
		return stripNewLinesPattern.ReplaceAll(typed, emptyBytes)
	default:
		return input
	}
}
