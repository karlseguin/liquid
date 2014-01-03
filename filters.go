package liquid

import (
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/filters"
)

// A filter factory creates a filter based on the supplied parameters
type FilterFactory func(parameters []core.Value) filters.Filter

// A map of filters. You can register custom filters by adding them to this
// map. Note however that the map isn't thread-safe.
var Filters = map[string]FilterFactory{
	"capitalize":     filters.CapitalizeFactory,
	"downcase":       filters.DowncaseFactory,
	"upcase":         filters.UpcaseFactory,
	"first":          filters.FirstFactory,
	"last":           filters.LastFactory,
	"join":           filters.JoinFactory,
	"debug":          filters.DebugFactory,
	"plus":           filters.PlusFactory,
	"minus":          filters.MinusFactory,
	"size":           filters.SizeFactory,
	"times":          filters.TimesFactory,
	"divideby":       filters.DivideByFactory,
	"prepend":        filters.PrependFactory,
	"append":         filters.AppendFactory,
	"strip_newlines": filters.StripNewLinesFactory,
	"replace":        filters.ReplaceFactory,
	"replace_first":  filters.ReplaceFirstFactory,
	"remove":         filters.RemoveFactory,
	"remove_first":   filters.RemoveFirstFactory,
	"newline_to_br":  filters.NewLineToBrFactory,
	"split":          filters.SplitFactory,
	"modulo":         filters.ModuloFactory,
	"truncate":       filters.TruncateFactory,
	"truncatewords":  filters.TruncateWordsFactory,
}
