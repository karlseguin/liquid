package liquid

type Literal struct {
	Value []byte
}

func literalExtractor(data []byte) (Token, error) {
	l := &Literal{
		Value: make([]byte, len(data)),
	}
	copy(l.Value, data)
	return l, nil
}

func (l *Literal) Render(data interface{}) []byte {
	return l.Value
}
