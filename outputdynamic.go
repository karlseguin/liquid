package liquid

import (
	"github.com/karlseguin/liquid/filters"
)

type OutputDynamic struct {
	Values [][]byte
	Filters []filters.Filter
}
