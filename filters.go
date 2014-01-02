package liquid

import (
	"github.com/karlseguin/liquid/filters"
)

// A filter factory creates a filter based on the supplied parameters
type FilterFactory func(parameters []string) filters.Filter

// A map of filters. You can register custom filters by adding them to this
// map. Note however that the map isn't thread-safe.
var Filters = map[string]FilterFactory{
	"capitalize": filters.CapitalizeFactory,
	"downcase":   filters.DowncaseFactory,
	"upcase":     filters.UpcaseFactory,
	"first":      filters.FirstFactory,
	"last":       filters.LastFactory,
	"join":       filters.JoinFactory,
	"debug":      filters.DebugFactory,
	"plus":       filters.PlusFactory,
}
