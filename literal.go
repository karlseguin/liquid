package liquid

import (
	"github.com/karlseguin/liquid/core"
)

type Literal struct {
	Value []byte
}

func newLiteral(data []byte) core.Token {
	l := &Literal{
		Value: make([]byte, len(data)),
	}
	copy(l.Value, data)
	return l
}

func (l *Literal) Render(data map[string]interface{}) []byte {
	return l.Value
}
