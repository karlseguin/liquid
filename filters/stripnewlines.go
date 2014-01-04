package filters

import (
	"github.com/karlseguin/liquid/core"
	"regexp"
)

var stripNewLines = &ReplacePattern{regexp.MustCompile("(\n|\r)"), ""}

func StripNewLinesFactory(parameters []core.Value) core.Filter {
	return stripNewLines.Replace
}
