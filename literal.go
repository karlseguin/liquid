package liquid

import (
	"github.com/karlseguin/liquid/core"
)

type Literal struct {
	Value []byte
}

func literalExtractor(data []byte) (core.Token, error) {
	l := &Literal{
		Value: make([]byte, len(data)),
	}
	copy(l.Value, data)
	return l, nil
}

func (l *Literal) Render(data map[string]interface{}) []byte {
	return l.Value
}
