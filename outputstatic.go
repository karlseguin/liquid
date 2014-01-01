package liquid

import (
	"github.com/karlseguin/liquid/filters"
)

type OutputStatic struct {
	Value []byte
	Filters []filters.Filter
}
