package filters

import (
	"github.com/karlseguin/liquid/core"
	"regexp"
)

var newLinesToBr = &ReplacePattern{regexp.MustCompile("(\n\r|\n|\r)"), "<br />\n"}

func NewLineToBrFactory(parameters []core.Value) core.Filter {
	return newLinesToBr.Replace
}
