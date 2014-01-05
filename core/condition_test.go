package core

import (
	"testing"
	"time"
)

func TestUnaryCondition(t *testing.T) {
	assertCondition(t, true, boolValue(true), Unary, nil)
	assertCondition(t, false, boolValue(false), Unary, nil)
	assertCondition(t, true, stringValue("spice"), Unary, nil)
	assertCondition(t, false, stringValue(""), Unary, nil)
	assertCondition(t, true, dynamicValue("string"), Unary, nil)
	assertCondition(t, true, dynamicValue("[]int"), Unary, nil)
	assertCondition(t, false, dynamicValue("notexists"), Unary, nil)
}

func TestEqualsConditionWithBalancedStrings(t *testing.T) {
	assertEqualsCondition(t, true, stringValue("abc"), stringValue("abc"))
	assertEqualsCondition(t, true, stringValue(""), stringValue(""))
	assertEqualsCondition(t, false, stringValue("abc"), stringValue("123"))
}

func TestEqualsConditionWithBalancedDynamicStrings(t *testing.T) {
	assertEqualsCondition(t, true, dynamicValue("doesnotexist"), dynamicValue("doesnotexist"))
	assertEqualsCondition(t, true, dynamicValue("string"), stringValue("astring"))
	assertEqualsCondition(t, false, dynamicValue("string"), stringValue("other"))
}

func TestEqualsConditionWithBalancedDynamicArrays(t *testing.T) {
	assertEqualsCondition(t, true, dynamicValue("[]int"), dynamicValue("[]int"))
	assertEqualsCondition(t, false, dynamicValue("[]int"), dynamicValue("[]int2"))
}

func TestEqualsConditionWithBalancedBools(t *testing.T) {
	assertEqualsCondition(t, true, boolValue(true), boolValue(true))
	assertEqualsCondition(t, false, boolValue(true), boolValue(false))
}

func TestEqualsConditionWithBalancedInt(t *testing.T) {
	assertEqualsCondition(t, true, intValue(3231), intValue(3231))
	assertEqualsCondition(t, false, intValue(3231), intValue(2993))
}

func TestEqualsConditionWithBalancedFloat(t *testing.T) {
	assertEqualsCondition(t, true, floatValue(11.33), floatValue(11.33))
	assertEqualsCondition(t, false, floatValue(11.2), floatValue(11.21))
}

func TestEqualWithUnbalancedInt(t *testing.T) {
	assertEqualsCondition(t, true, intValue(123), stringValue("123"))
	assertEqualsCondition(t, false, intValue(123), stringValue("1a23"))
	assertEqualsCondition(t, true, intValue(123), floatValue(123.0))
	assertEqualsCondition(t, false, intValue(123), floatValue(123.1))
}

func TestEqualWithUnbalancedFloats(t *testing.T) {
	assertEqualsCondition(t, true, floatValue(123.0), stringValue("123"))
	assertEqualsCondition(t, true, floatValue(123.0), intValue(123))
	assertEqualsCondition(t, false, floatValue(123.0), stringValue("123.1"))
}

func TestEqualsWithTime(t *testing.T) {
	assertEqualsCondition(t, true, dynamicValue("now"), stringValue("today"))
	assertEqualsCondition(t, false, dynamicValue("yesterday"), stringValue("today"))
}

func TestConditionGroupWithOneCondition(t *testing.T) {
	assertConditionGroup(t, true, trueCondition())
	assertConditionGroup(t, false, falseCondition())
}

func TestConditionGroupWithTwoOrCondition(t *testing.T) {
	assertConditionGroup(t, true, trueCondition(), OR, trueCondition())
	assertConditionGroup(t, true, trueCondition(), OR, falseCondition())
	assertConditionGroup(t, true, falseCondition(), OR, trueCondition())
	assertConditionGroup(t, false, falseCondition(), OR, falseCondition())
}

func TestConditionGroupWithTwoAndCondition(t *testing.T) {
	assertConditionGroup(t, true, trueCondition(), AND, trueCondition())
	assertConditionGroup(t, false, trueCondition(), AND, falseCondition())
	assertConditionGroup(t, false, falseCondition(), AND, trueCondition())
	assertConditionGroup(t, false, falseCondition(), AND, falseCondition())
}

func TestConditionGroupWithMultipleConditions(t *testing.T) {
	assertConditionGroup(t, true, trueCondition(), OR, trueCondition(), AND, falseCondition())
	assertConditionGroup(t, true, trueCondition(), AND, trueCondition(), OR, trueCondition())
	assertConditionGroup(t, false, falseCondition(), OR, trueCondition(), AND, falseCondition())
	assertConditionGroup(t, false, falseCondition(), OR, trueCondition(), AND, falseCondition(), OR, falseCondition())
	assertConditionGroup(t, true, falseCondition(), OR, trueCondition(), AND, falseCondition(), OR, trueCondition())
}

func assertEqualsCondition(t *testing.T, expected bool, left, right Value) {
	assertCondition(t, expected, left, Equals, right)
	assertCondition(t, !expected, left, NotEquals, right)
	if expected {
		assertCondition(t, false, left, LessThan, right)
		assertCondition(t, false, left, GreaterThan, right)
		assertCondition(t, true, left, LessThanOrEqual, right)
		assertCondition(t, true, left, GreaterThanOrEqual, right)
	}
}

func assertCondition(t *testing.T, expected bool, left Value, op ComparisonOperator, right Value) {
	data := map[string]interface{}{
		"[]int":     []int{1, 2, 3},
		"[]int2":    []int{2, 3, 1},
		"string":    "astring",
		"now":       time.Now(),
		"yesterday": time.Now().Add(time.Hour * -24),
	}
	c := &Condition{left, op, right}
	actual := c.IsTrue(data)
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func assertConditionGroup(t *testing.T, expected bool, data ...interface{}) {
	l := len(data)
	group := &ConditionGroup{
		joins:      make([]ConditionGroupJoin, 0, l/2),
		conditions: make([]*Condition, 0, l-l/2),
	}
	for i := 0; i < l; i += 2 {
		group.conditions = append(group.conditions, data[i].(*Condition))
		if i+1 < l {
			group.joins = append(group.joins, data[i+1].(ConditionGroupJoin))
		}
	}

	actual := group.IsTrue(nil)
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func boolValue(b bool) Value {
	return &StaticBoolValue{b}
}

func stringValue(s string) Value {
	return &StaticStringValue{[]byte(s)}
}

func intValue(n int) Value {
	return &StaticIntValue{n}
}

func floatValue(f float64) Value {
	return &StaticFloatValue{f}
}

func dynamicValue(s string) Value {
	return &DynamicValue{[]string{s}}
}

func trueCondition() *Condition {
	return &Condition{boolValue(true), Unary, nil}
}

func falseCondition() *Condition {
	return &Condition{boolValue(false), Unary, nil}
}
