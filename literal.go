package liquid

import (
	"github.com/karlseguin/liquid/core"
)

type Literal struct {
	Value []byte
}

// Creates a literal (just plain text)
func newLiteral(data []byte) core.Code {
	return &Literal{Value: data}
}

func (l *Literal) Render(data map[string]interface{}) []byte {
	return l.Value
}
