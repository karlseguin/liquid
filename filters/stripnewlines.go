package filters

import (
	"github.com/karlseguin/liquid/core"
	"regexp"
)

var stripNewLinesPattern = regexp.MustCompile("(\n|\r)")
var emptyBytes = []byte("")

func StripNewLinesFactory(parameters []core.Value) Filter {
	return StripNewLines
}

func StripNewLines(input interface{}, data map[string]interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return stripNewLinesPattern.ReplaceAllString(typed, "")
	case []byte:
		return stripNewLinesPattern.ReplaceAll(typed, emptyBytes)
	default:
		return input
	}
}
