package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	EmptyBytes = []byte{}
)

type MarkupType int

const (
	OutputMarkup MarkupType = iota
	TagMarkup
	NoMarkup
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

func (p *Parser) ToMarkup() ([]byte, MarkupType) {
	start := p.Position
	markupType := NoMarkup
	for ; p.Position < p.Len; p.Position++ {
		if p.SkipUntil('{') != '{' {
			break
		}
		next := p.Peek()
		if next == '{' {
			markupType = OutputMarkup
			break
		}
		if next == '%' {
			markupType = TagMarkup
			break
		}
	}
	pre := EmptyBytes
	if p.Position > start {
		pre = p.Data[start:p.Position]
	}
	p.Commit()
	return pre, markupType
}

func (p *Parser) SkipPastTag() {
	for p.HasMore() {
		p.SkipUntil('}')
		p.Forward()
		if p.Data[p.Position-2] == '%' {
			return
		}
	}
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

func (p *Parser) ReadValue() (Value, error) {
	current := p.SkipSpaces()
	if current == 0 || current == '}' || current == '|' || current == ':' || current == '%' {
		return nil, nil
	}
	if current == '\'' || current == '"' {
		return p.ReadStaticStringValue(current)
	}
	if current == '-' || (current >= '0' && current <= '9') {
		return p.ReadStaticNumericValue()
	}
	return p.ReadDynamicValues()
}

func (p *Parser) ReadStaticStringValue(delimiter byte) (Value, error) {
	p.Forward() //consume the opening '
	escaped := 0
	found := false
	start := p.Position
	for {
		current := p.SkipUntil(delimiter)
		if p.Data[p.Position-1] != '\\' {
			if current == delimiter {
				p.Forward()
				found = true
			}
			break
		}
		p.Forward()
		escaped++
	}

	if found == false {
		return nil, p.Error("Invalid value, a single quote might be missing", start)
	}
	p.Commit()
	var data []byte
	if escaped > 0 {
		data = unescape(p.Data[start:p.Position-1], escaped)
	} else {
		data = detatch(p.Data[start : p.Position-1])
	}
	return &StaticStringValue{data}, nil
}

func (p *Parser) ReadStaticNumericValue() (Value, error) {
	start := p.Position
	name := p.ReadName()
	if len(name) == 0 {
		return nil, p.Error("Was expecting a value, got nothing", start)
	}
	if i, err := strconv.Atoi(name); err == nil {
		return &StaticIntValue{i}, nil
	}
	if f, err := strconv.ParseFloat(name, 64); err == nil {
		return &StaticFloatValue{f}, nil
	}
	return nil, p.Error("Invalid value. This is either an invalid number, variable name, or maybe you're missing a quote", start)
}

func (p *Parser) ReadDynamicValues() (Value, error) {
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
	return &DynamicValue{TrimStrings(values)}, nil
}

func (p *Parser) ReadName() string {
	var name string
	p.SkipSpaces()
	marker := p.Position
	for ; p.Position < p.Len; p.Position++ {
		current := p.Current()
		if current == ' ' || current == '|' || current == '}' || current == '%' || current == ':' || current == ',' {
			name = string(p.Data[marker:p.Position])
			break
		}
	}
	p.Commit()
	return name
}

func (p *Parser) ReadParameters() ([]Value, error) {
	values := make([]Value, 0, 3)
	for {
		value, err := p.ReadValue()
		if err != nil {
			return nil, err
		}
		if value == nil {
			break
		}
		values = append(values, value)
		if p.SkipSpaces() == ',' {
			p.Forward()
			continue
		}
		break
	}
	p.Commit()
	return TrimValues(values), nil
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
	p.ForwardBy(1)
}

func (p *Parser) ForwardBy(count int) {
	p.Position += count
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
	start = start - 5
	if start < 0 {
		start = 0
	}
	end := start + 20
	if end > p.Len {
		end = p.Len
	}
	return p.Data[start:end]
}

func (p *Parser) Out() {
	fmt.Println(string(p.Data[p.Position:]))
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
