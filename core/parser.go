package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	EmptyBytes = []byte{}
	Whitespace = []byte(" ")
)

type MarkupType int

const (
	OutputMarkup MarkupType = iota
	TagMarkup
	NoMarkup
)

type Parser struct {
	Position int
	Data     []byte
	Len      int
	End      int
	Line     int
}

func NewParser(data []byte) *Parser {
	parser := &Parser{
		Position: 0,
		Data:     data,
		Len:      len(data),
		Line:     1,
	}
	parser.End = parser.Len - 1
	return parser
}

func (p *Parser) ToMarkup(preserveWhitespace bool) ([]byte, MarkupType) {
	isWhitespace := true
	start := p.Position
	markupType := NoMarkup
	for ; p.Position < p.Len; p.Position++ {
		current := p.Current()
		if current == '{' {
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
		if isWhitespace && current != ' ' && current != '\r' && current != '\n' {
			isWhitespace = false
		}
	}
	pre := EmptyBytes
	if p.Position > start {
		if isWhitespace && !preserveWhitespace {
			pre = Whitespace
		} else {
			pre = detach(p.Data[start:p.Position])
		}
	}
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

func (p *Parser) SkipPastOutput() {
	for p.HasMore() {
		p.SkipUntil('}')
		p.Forward()
		if p.Data[p.Position-2] == '}' {
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
	return
}

func (p *Parser) ReadValue() (Value, error) {
	current := p.SkipSpaces()
	if IsTokenEnd(current) {
		return nil, nil
	}
	if current == '\'' || current == '"' {
		return p.ReadStaticStringValue(current)
	}
	if current == '-' || (current >= '0' && current <= '9') {
		return p.ReadStaticNumericValue()
	}
	if current == '(' {
		return p.ReadRangeValue()
	}
	if b, ok := p.ReadStaticBoolValue(); ok {
		return b, nil
	}
	if e, ok := p.ReadStaticEmptyValue(); ok {
		return e, nil
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
		return nil, p.Error("Invalid value, a single quote might be missing")
	}
	var data string
	if escaped > 0 {
		data = string(unescape(p.Data[start:p.Position-1], escaped))
	} else {
		data = string(p.Data[start : p.Position-1])
	}
	return &StaticStringValue{data}, nil
}

func (p *Parser) ReadRangeValue() (Value, error) {
	p.Forward() //consume the opening (
	from, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	if from == nil {
		return nil, p.Error("Invalid range, should start with a number of variable")
	}
	if p.SkipUntil('.') != '.' || p.Left() < 2 || p.Data[p.Position+1] != '.' {
		return nil, p.Error("Invalid range, expecting '..' between two values")
	}
	p.ForwardBy(2)
	to, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	if to == nil {
		return nil, p.Error("Invalid range, should end with a number of variable")
	}

	if p.SkipUntil(')') != ')' {
		return nil, p.Error("Invalid range, expecting a close ')'")
	}
	p.Forward()
	return &RangeValue{from, to}, nil
}

func (p *Parser) ReadStaticNumericValue() (Value, error) {
	hasDecimal := false
	var value string
	p.SkipSpaces()
	marker := p.Position
	for ; p.Position < p.Len; p.Position++ {
		current := p.Current()
		if !IsTokenEnd(current) {
			continue
		}
		if current == '.' {
			next := p.Peek()
			if next <= '9' && next >= '0' && !hasDecimal {
				hasDecimal = true
				continue
			}
		}
		value = string(p.Data[marker:p.Position])
		break
	}
	if len(value) == 0 {
		return nil, p.Error("Was expecting a value, got nothing")
	}
	if i, err := strconv.Atoi(value); err == nil {
		return &StaticIntValue{i}, nil
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return &StaticFloatValue{f}, nil
	}
	return nil, p.Error("Invalid value. It looked like a number but could not be convert")
}

func (p *Parser) ReadStaticBoolValue() (Value, bool) {
	start := p.Position
	left := p.Left()
	if left > 4 {
		if p.Data[start] == 't' && p.Data[start+1] == 'r' && p.Data[start+2] == 'u' && p.Data[start+3] == 'e' && IsTokenEnd(p.Data[start+4]) {
			p.ForwardBy(4)
			return &StaticBoolValue{true}, true
		}
	}
	if left > 5 {
		if p.Data[start] == 'f' && p.Data[start+1] == 'a' && p.Data[start+2] == 'l' && p.Data[start+3] == 's' && p.Data[start+4] == 'e' && IsTokenEnd(p.Data[start+5]) {
			p.ForwardBy(5)
			return &StaticBoolValue{false}, true
		}
	}
	return nil, false
}

func (p *Parser) ReadStaticEmptyValue() (Value, bool) {
	start := p.Position
	if p.Left() > 5 {
		if p.Data[start] == 'e' && p.Data[start+1] == 'm' && p.Data[start+2] == 'p' && p.Data[start+3] == 't' && p.Data[start+4] == 'y' && IsTokenEnd(p.Data[start+5]) {
			p.ForwardBy(5)
			return &StaticEmptyValue{}, true
		}
	}
	return nil, false
}

func (p *Parser) ReadDynamicValues() (Value, error) {
	values := make([]string, 0, 5)
	marker := p.Position
	for ; p.Position < p.Len; p.Position++ {
		current := p.Current()
		if current == '.' {
			values = append(values, strings.ToLower(string(p.Data[marker:p.Position])))
			marker = p.Position + 1
		} else if IsTokenEnd(current) {
			values = append(values, strings.ToLower(string(p.Data[marker:p.Position])))
			break
		}
	}
	return NewDynamicValue(TrimStrings(values)), nil
}

func (p *Parser) ReadName() string {
	var name string
	p.SkipSpaces()
	marker := p.Position
	for ; p.Position < p.Len; p.Position++ {
		current := p.Current()
		if IsTokenEnd(current) {
			name = string(p.Data[marker:p.Position])
			break
		}
	}
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
	return TrimValues(values), nil
}

func (p *Parser) ReadFilters() ([]Filter, error) {
	var filters []Filter
	if p.SkipSpaces() != '|' {
		return filters, nil
	}
	p.Forward()
	filters = make([]Filter, 0, 5)
	for name := p.ReadName(); name != ""; name = p.ReadName() {
		factory, exists := FilterLookup[name]
		if exists == false {
			return nil, p.Error(fmt.Sprintf("Unknown filter %q", name))
		}
		var parameters []Value
		if p.SkipSpaces() == ':' {
			p.Forward()
			var err error
			parameters, err = p.ReadParameters()
			if err != nil {
				return nil, err
			}
		}
		filters = append(filters, factory(parameters))
		if p.SkipSpaces() == '|' {
			p.Forward()
			continue
		}
		break
	}
	return filters, nil
}

func (p *Parser) ReadConditionGroup() (Verifiable, error) {
	group := &ConditionGroup{make([]*Condition, 0, 2), make([]LogicalOperator, 0, 1), false}
	for {
		left, err := p.ReadValue()
		if err != nil {
			return nil, err
		}
		if left == nil {
			return nil, p.Error("Invalid of missing left value in condition")
		}

		if p.SkipSpaces() == '%' {
			group.conditions = append(group.conditions, &Condition{left, Unary, nil})
			break
		}
		operator := p.ReadComparisonOperator()
		if operator == UnknownComparator {
			logical := p.ReadLogicalOperator()
			if logical == UnknownLogical {
				return nil, p.Error("Invalid or missing operator (should be ==, !=, >, <, >=, <= or contains)")
			}
			group.conditions = append(group.conditions, &Condition{left, Unary, nil})
			group.joins = append(group.joins, logical)
			continue
		}

		right, err := p.ReadValue()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.Error("Invalid of missing right value in condition")
		}
		group.conditions = append(group.conditions, &Condition{left, operator, right})

		if p.SkipSpaces() == '%' {
			break
		}

		logical := p.ReadLogicalOperator()
		if logical == UnknownLogical {
			return nil, p.Error("Invalid condition. Expecting 'and', 'or' or end of tag")
		}
		group.joins = append(group.joins, logical)
	}
	return group, nil
}

func (p *Parser) ReadComparisonOperator() ComparisonOperator {
	left := p.Left()
	start := p.Position
	if left > 1 {
		current := p.Data[start]
		next := p.Data[start+1]
		if current == '=' && next == '=' && IsTokenEnd(p.Data[start+2]) {
			p.ForwardBy(2)
			return Equals
		}
		if current == '!' && next == '=' && IsTokenEnd(p.Data[start+2]) {
			p.ForwardBy(2)
			return NotEquals
		}
		if current == '>' && IsTokenEnd(next) {
			p.ForwardBy(1)
			return GreaterThan
		}
		if current == '<' && IsTokenEnd(next) {
			p.ForwardBy(1)
			return LessThan
		}
		if current == '>' && next == '=' && IsTokenEnd(p.Data[start+2]) {
			p.ForwardBy(2)
			return GreaterThanOrEqual
		}
		if current == '<' && next == '=' && IsTokenEnd(p.Data[start+2]) {
			p.ForwardBy(1)
			return LessThanOrEqual
		}
		if left > 7 {
			if current == 'c' && next == 'o' && p.Data[start+2] == 'n' && p.Data[start+3] == 't' && p.Data[start+4] == 'a' && p.Data[start+5] == 'i' && p.Data[start+6] == 'n' && p.Data[start+7] == 's' && IsTokenEnd(p.Data[start+8]) {
				p.ForwardBy(8)
				return Contains
			}
		}
	}
	return UnknownComparator
}

func (p *Parser) ReadLogicalOperator() LogicalOperator {
	left := p.Left()
	start := p.Position
	if left > 2 {
		current := p.Data[start]
		next := p.Data[start+1]
		if current == 'o' && next == 'r' && IsTokenEnd(p.Data[start+2]) {
			p.ForwardBy(2)
			return OR
		}
		if current == 'a' && next == 'n' && p.Data[start+2] == 'd' && IsTokenEnd(p.Data[start+3]) {
			p.ForwardBy(3)
			return AND
		}
	}
	return UnknownLogical
}

func (p *Parser) ReadPartialCondition() (Completable, error) {
	group := &ConditionGroup{make([]*Condition, 0, 2), make([]LogicalOperator, 0, 1), false}
	for {
		value, err := p.ReadValue()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, p.Error("Invalid of missing value for condition")
		}
		group.conditions = append(group.conditions, &Condition{value, UnknownComparator, nil})

		if p.SkipSpaces() == '%' {
			break
		}
		logical := p.ReadLogicalOperator()
		if logical == UnknownLogical {
			return nil, p.Error("Invalid condition. Expecting 'and', 'or' or end of tag")
		}
		group.joins = append(group.joins, logical)
	}
	return group, nil
}

func (p *Parser) Peek() byte {
	if p.Position == p.End {
		return 0
	}
	return p.Data[p.Position+1]
}

func (p *Parser) Left() int {
	return p.Len - p.Position
}
func (p *Parser) HasMore() bool {
	return p.Position < p.End
}

func (p *Parser) Current() byte {
	current := p.Data[p.Position]
	if current == '\n' {
		p.Line++
	}
	return current
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
	return
}

func (p *Parser) Error(s string) error {
	end := p.Position
	if end == p.Len {
		end--
	} else {
		for ; end < p.Len; end++ {
			if p.Data[end] == '}' {
				if end < p.End && p.Data[end+1] == '}' {
					end++
				}
				break
			}
		}
	}

	start := end - 25
	if start < 0 {
		start = 0
	}
	line := p.Line
	for i := end; i > 0; i-- {
		if p.Data[i] == '\n' {
			line--
		} else if p.Data[i] == '{' || p.Data[i] == '%' && p.Data[i-1] == '{' {
			start = i - 1
			break
		}
	}
	return errors.New(fmt.Sprintf("%s (%q - line %d)", s, string(p.Data[start:end+1]), line))
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

func IsTokenEnd(b byte) bool {
	return b == ' ' || b == '|' || b == '}' || b == '%' || b == ':' || b == ',' || b == ')' || b == '.' || b == 0
}

func detach(data []byte) []byte {
	detached := make([]byte, len(data))
	copy(detached, data)
	return detached
}
