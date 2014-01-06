package core

import (
	"bytes"
	"reflect"
	"strings"
	"time"
)

type LogicalOperator int
type ComparisonOperator int
type Type int

const (
	OR LogicalOperator = iota
	AND
	UnknownLogical

	Equals ComparisonOperator = iota
	NotEquals
	LessThan
	GreaterThan
	LessThanOrEqual
	GreaterThanOrEqual
	Contains
	NotContains
	Unary
	NotUnary
	UnknownComparator

	String Type = iota
	Nil
	Int
	Int64
	Uint
	Float64
	Complex128
	Bool
	Time
	Today
	Array
	Unknown
)

var KindToType = map[reflect.Kind]Type{
	reflect.String:     String,
	reflect.Int:        Int,
	reflect.Int64:      Int64,
	reflect.Uint:       Uint,
	reflect.Float64:    Float64,
	reflect.Complex128: Complex128,
	reflect.Bool:       Bool,
	reflect.Array:      Array,
	reflect.Slice:      Array,
	reflect.Map:        Array,
}

var TypeOperations = map[Type]map[ComparisonOperator]ConditionResolver{
	String: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(string) == right.(string) },
		LessThan: func(left, right interface{}) bool { return left.(string) < right.(string) },
	},
	Nil: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left == nil && right == nil },
		LessThan: func(left, right interface{}) bool { return false },
	},
	Int: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(int) == right.(int) },
		LessThan: func(left, right interface{}) bool { return left.(int) < right.(int) },
	},
	Int64: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(int64) == right.(int64) },
		LessThan: func(left, right interface{}) bool { return left.(int64) < right.(int64) },
	},
	Uint: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(uint) == right.(uint) },
		LessThan: func(left, right interface{}) bool { return left.(uint) < right.(uint) },
	},
	Float64: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(float64) == right.(float64) },
		LessThan: func(left, right interface{}) bool { return left.(float64) < right.(float64) },
	},
	Complex128: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(complex128) == right.(complex128) },
		LessThan: func(left, right interface{}) bool { return false },
	},
	Bool: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(bool) == right.(bool) },
		LessThan: func(left, right interface{}) bool { return false },
	},
	Time: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return left.(time.Time).Unix() == right.(time.Time).Unix() },
		LessThan: func(left, right interface{}) bool { return left.(time.Time).Unix() < right.(time.Time).Unix() },
	},
	Today: map[ComparisonOperator]ConditionResolver{
		Equals: func(left, right interface{}) bool {
			l, r := left.(time.Time), right.(time.Time)
			return l.YearDay() == r.YearDay() && l.Year() == r.Year()
		},
		LessThan: func(left, right interface{}) bool {
			l, r := left.(time.Time), right.(time.Time)
			if l.Year() > r.Year() {
				return false
			}
			if l.Year() < r.Year() {
				return true
			}
			return l.YearDay() < r.YearDay()
		},
	},
	Array: map[ComparisonOperator]ConditionResolver{
		Equals:   func(left, right interface{}) bool { return reflect.DeepEqual(left, right) },
		LessThan: func(left, right interface{}) bool { return reflect.ValueOf(left).Len() < reflect.ValueOf(right).Len() },
	},
}

// Resolves a condition
type ConditionResolver func(left, right interface{}) bool

var ConditionLookup = map[ComparisonOperator]ConditionResolver{
	Unary:              UnaryComparison,
	NotUnary:           NotUnaryComparison,
	Equals:             EqualsComparison,
	NotEquals:          NotEqualsComparison,
	LessThan:           LessThanComparison,
	GreaterThan:        GreaterThanComparison,
	LessThanOrEqual:    LessThanOrEqualComparison,
	GreaterThanOrEqual: GreaterThanOrEqualComparison,
	Contains:           ContainsComparison,
}

type Completable interface {
	Complete(value Value, operator ComparisonOperator)
	Verifiable
}

type Verifiable interface {
	IsTrue(data map[string]interface{}) bool
	Inverse()
}

// represents a group of conditions
type ConditionGroup struct {
	conditions []*Condition
	joins      []LogicalOperator
	inverse    bool
}

func (g *ConditionGroup) Inverse() {
	g.inverse = true
}

func (g ConditionGroup) Complete(value Value, operator ComparisonOperator) {
	for _, condition := range g.conditions {
		condition.right = value
		condition.operator = operator
	}
}

func (g *ConditionGroup) IsTrue(data map[string]interface{}) bool {
	l := len(g.conditions) - 1
	if l == 0 {
		return g.realReturn(g.conditions[0].IsTrue(data))
	}

	for i := 0; i <= l; i++ {
		if g.conditions[i].IsTrue(data) {
			if i == l || g.joins[i] == OR {
				return g.realReturn(true)
			}
		} else if i != l && g.joins[i] == AND {
			for ; i < l; i++ {
				if g.joins[i] == OR {
					break
				}
			}
		}
	}
	return g.realReturn(false)
}

func (g *ConditionGroup) realReturn(b bool) bool {
	if g.inverse {
		return !b
	}
	return b
}

type TrueCondition struct {
	inverse bool
}

func (t *TrueCondition) IsTrue(data map[string]interface{}) bool {
	if t.inverse {
		return false
	}
	return true
}

func (t *TrueCondition) Inverse() {
	t.inverse = true
}

// represents a conditions (such as x == y)
type Condition struct {
	left     Value
	operator ComparisonOperator
	right    Value
}

func (c *Condition) IsTrue(data map[string]interface{}) bool {
	left := c.left.ResolveWithNil(data)
	var right interface{}
	if c.right != nil {
		right = c.right.ResolveWithNil(data)
	}
	return ConditionLookup[c.operator](left, right)
}

func UnaryComparison(left, right interface{}) bool {
	if left == nil {
		return false
	}
	switch typed := left.(type) {
	case bool:
		return typed
	case string:
		return len(typed) > 0
	case []byte:
		return len(typed) > 0
	}
	return true
}

func NotUnaryComparison(left, right interface{}) bool {
	return !UnaryComparison(left, right)
}

func EqualsComparison(left, right interface{}) bool {
	if s, ok := right.(string); ok && s == "liquid:empty" {
		if n, ok := ToLength(left); ok {
			return n == 0
		}
		return false
	}

	var t Type
	if left, right, t = convertToSameType(left, right); t == Unknown {
		return false
	}
	return TypeOperations[t][Equals](left, right)
}

func NotEqualsComparison(left, right interface{}) bool {
	return !EqualsComparison(left, right)
}

func LessThanComparison(left, right interface{}) bool {
	var t Type
	if left, right, t = convertToSameType(left, right); t == Unknown {
		return false
	}
	return TypeOperations[t][LessThan](left, right)
}

func LessThanOrEqualComparison(left, right interface{}) bool {
	var t Type
	if left, right, t = convertToSameType(left, right); t == Unknown {
		return false
	}
	return TypeOperations[t][Equals](left, right) || TypeOperations[t][LessThan](left, right)
}

func GreaterThanComparison(left, right interface{}) bool {
	var t Type
	if left, right, t = convertToSameType(left, right); t == Unknown {
		return false
	}
	return !TypeOperations[t][Equals](left, right) && !TypeOperations[t][LessThan](left, right)
}

func GreaterThanOrEqualComparison(left, right interface{}) bool {
	var t Type
	if left, right, t = convertToSameType(left, right); t == Unknown {
		return false
	}
	return !TypeOperations[t][LessThan](left, right)
}

// I think most of this sucks
func ContainsComparison(left, right interface{}) bool {
	if s, ok := left.(string); ok {
		return strings.Contains(s, ToString(right))
	}
	if b, ok := left.([]byte); ok {
		return bytes.Contains(b, ToBytes(right))
	}

	if strs, ok := left.([]string); ok {
		needle := ToString(right)
		for i, l := 0, len(strs); i < l; i++ {
			if strs[i] == needle {
				return true
			}
		}
		return false
	}

	if n, ok := left.([]int); ok {
		needle, ok := ToInt(right)
		if ok == false {
			return false
		}
		for i, l := 0, len(n); i < l; i++ {
			if n[i] == needle {
				return true
			}
		}
		return false
	}

	value := reflect.ValueOf(left)
	kind := value.Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		l := value.Len()
		if l == 0 {
			return false
		}
		for i := 0; i < l; i++ {
			if EqualsComparison(value.Index(i).Interface(), right) {
				return true
			}
		}
		return false
	}
	if kind == reflect.Map {
		if value.Len() == 0 {
			return false
		}
		if b, ok := right.([]byte); ok {
			right = string(b)
		}
		rightValue := reflect.ValueOf(right)
		if rightValue.Type() == value.Type().Key() {
			return value.MapIndex(rightValue).IsValid()
		}
		return false
	}
	return false
}

func NotContainsComparison(left, right interface{}) bool {
	return !ContainsComparison(left, right)
}

func convertToSameType(left, right interface{}) (interface{}, interface{}, Type) {
	//rely on the above code to handle this properly
	if left == nil || right == nil {
		return left, right, Nil
	}

	if s, ok := left.(string); ok {
		return convertStringsToSameType(s, right)
	} else if s, ok := right.(string); ok {
		return convertStringsToSameType(s, left)
	}
	if b, ok := left.([]byte); ok {
		return convertStringsToSameType(string(b), right)
	} else if b, ok := right.([]byte); ok {
		return convertStringsToSameType(string(b), left)
	}

	leftValue, rightValue := reflect.ValueOf(left), reflect.ValueOf(right)
	leftKind, rightKind := leftValue.Kind(), rightValue.Kind()
	if leftKind == rightKind {
		if t, ok := KindToType[leftKind]; ok {
			return left, right, t
		}
	}
	if left, right, t := convertNumbersToSameType(leftValue, leftKind, rightValue, rightKind); t != Unknown {
		return left, right, t
	}
	return left, right, Unknown
}

func convertStringsToSameType(a string, b interface{}) (interface{}, interface{}, Type) {
	if a == "today" {
		if t, ok := b.(time.Time); ok {
			return Now(), t, Today
		}
	} else if a == "now" {
		if t, ok := b.(time.Time); ok {
			return Now(), t, Time
		}
	}
	return a, ToString(b), String
}

func convertNumbersToSameType(leftValue reflect.Value, leftKind reflect.Kind, rightValue reflect.Value, rightKind reflect.Kind) (interface{}, interface{}, Type) {
	if isInt(leftKind) {
		if isInt(rightKind) {
			return leftValue.Int(), rightValue.Int(), Int64
		} else if isFloat(rightKind) {
			return float64(leftValue.Int()), rightValue.Float(), Float64
		}
	} else if isFloat(leftKind) {
		if isInt(rightKind) {
			return leftValue.Float(), float64(rightValue.Int()), Float64
		} else if isFloat(rightKind) {
			return leftValue.Float(), rightValue.Float(), Float64
		}
	} else if isComplex(leftKind) && isComplex(rightKind) {
		return leftValue.Complex(), rightValue.Complex(), Complex128
	}
	return nil, nil, Unknown
}

func isInt(kind reflect.Kind) bool {
	return kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64
}

func isFloat(kind reflect.Kind) bool {
	return kind == reflect.Float64 || kind == reflect.Float32
}

func isComplex(kind reflect.Kind) bool {
	return kind == reflect.Complex128 || kind == reflect.Complex64
}
