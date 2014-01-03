package core

import (
	"errors"
	"fmt"
	"strings"
)

const (
	OutputMarkup = 1
	TagMarkup    = 2
)

var (
	EmptyBytes = []byte{}
)

type Parser struct {
	Position       int
	Data           []byte
	Len            int
	End            int
	UncommitedLine int
	Line           int
}

func NewParser(data []byte) *Parser {
	parser := &Parser{
		Position:       0,
		Data:           data,
		Len:            len(data),
		Line:           1,
		UncommitedLine: 1,
	}
	parser.End = parser.Len - 1
	return parser
}

func (p *Parser) ToMarkup() ([]byte, int) {
	start := p.Position
	markupType := 0
	for ; p.Position < p.Len; p.Position++ {
		if p.Current() != '{' {
			continue
		}
		next := p.Peek()
		if next == '{' || next == '%' {
			p.Forward()
			markupType = OutputMarkup
			if next == '%' {
				markupType = TagMarkup
			}
			break
		}
	}
	pre := EmptyBytes
	if p.Position > start {
		if markupType != 0 {
			pre = p.Data[start : p.Position-2]
		} else {
			pre = p.Data[start:p.Position]
		}
	}
	p.Commit()
	return pre, markupType
}

func (p *Parser) SkipSpaces() (current byte) {
	for ; p.Position < p.Len; p.Position++ {
		c := p.Current()
		if c != ' ' {
			current = c
			break
		}
	}
	p.Commit()
	return
}

func (p *Parser) ReadValue() (interface{}, bool, error) {
	current := p.SkipSpaces()
	if current == 0 || current == '}' {
		return EmptyBytes, true, nil
	}
	if current == '\'' {
		p.Forward()
		static, err := p.ReadStaticValue()
		if err != nil {
			return nil, true, err
		}
		return static, true, nil
	}
	return p.ReadDynamicValues(), false, nil
}

func (p *Parser) ReadStaticValue() ([]byte, error) {
	escaped := 0
	start := p.Position
	for {
		p.SkipUntil('\'')
		if p.Data[p.Position-1] != '\\' {
			break
		}
		p.Forward()
		escaped++
	}

	if p.HasMore() == false {
		return EmptyBytes, p.Error("Invalid output value, a single quote might be missing", start)
	}
	p.Forward()
	p.Commit()
	if escaped > 0 {
		return unescape(p.Data[start:p.Position-1], escaped), nil
	}
	return detatch(p.Data[start : p.Position-1]), nil
}

func (p *Parser) ReadDynamicValues() []string {
	values := make([]string, 0, 5)
	marker := p.Position
	for ; p.Position < p.Len; p.Position++ {
		current := p.Current()
		if current == '.' {
			values = append(values, strings.ToLower(string(p.Data[marker:p.Position])))
			marker = p.Position + 1
		} else if current == '}' || current == ' ' || current == '|' {
			values = append(values, strings.ToLower(string(p.Data[marker:p.Position])))
			break
		}
	}
	p.Commit()
	return TrimStrings(values)
}

func (p *Parser) Peek() byte {
	if p.Position == p.End {
		return 0
	}
	return p.Data[p.Position+1]
}

func (p *Parser) HasMore() bool {
	return p.Position < p.End
}

func (p *Parser) Current() byte {
	current := p.Data[p.Position]
	if current == '\n' {
		p.UncommitedLine++
	}
	return current
}

func (p *Parser) Commit() {
	p.Line = p.UncommitedLine
}

func (p *Parser) Forward() {
	p.Position++
}

func (p *Parser) SkipUntil(b byte) (current byte) {
	for ; p.Position < p.Len; p.Position++ {
		if p.Current() == b {
			current = b
			break
		}
	}
	p.Commit()
	return
}

func (p *Parser) Error(s string, start int) error {
	return errors.New(fmt.Sprintf("%s (%q - line %d)", s, p.Snapshot(start), p.Line))
}

func (p *Parser) Snapshot(start int) []byte {
	start = start - 10
	if start < 0 {
		start = 0
	}
	end := start + 30
	if end > p.Len {
		end = p.Len
	}
	return p.Data[start:end]
}

func unescape(data []byte, escaped int) []byte {
	value := make([]byte, len(data)-escaped)
	i := 0
	found := 0
	position := 0
	for l := len(data) - 1; i < l; i++ {
		b := data[i]
		if b == '\\' && data[i+1] == '\'' {
			value[position] = '\''
			found++
			i++
			if found == escaped {
				break
			}
		} else {
			value[position] = b
		}
		position++
	}
	copy(value[position:], data[i:])
	return value
}

func detatch(data []byte) []byte {
	detached := make([]byte, len(data))
	copy(detached, data)
	return detached
}
