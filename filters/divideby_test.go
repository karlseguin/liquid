package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestDivideByAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{intValue(5)})
	spec.Expect(filter(43, nil).(float64)).ToEqual(8.6)
}

func TestDivideByAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{intValue(2)})
	spec.Expect(filter(43.3, nil).(float64)).ToEqual(21.65)
}

func TestDivideByAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{intValue(7)})
	spec.Expect(filter("33", nil).(float64)).ToEqual(4.714285714285714)
}

func TestDivideByAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{intValue(7)})
	spec.Expect(filter([]byte("34"), nil).(float64)).ToEqual(4.857142857142857)
}

func TestDivideByAnIntToAStringAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{intValue(7)})
	spec.Expect(filter("abc", nil).(string)).ToEqual("abc")
}

func TestDivideByAnIntToBytesAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{intValue(8)})
	spec.Expect(filter([]byte("abb"), nil).(string)).ToEqual("abb")
}

func TestDivideByAFloatToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{floatValue(1.10)})
	spec.Expect(filter(43, nil).(float64)).ToEqual(39.090909090909086)
}

func TestDivideByAFloatToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{floatValue(5.3)})
	spec.Expect(filter(43.3, nil).(float64)).ToEqual(8.169811320754716)
}

func TestDivideByAFloatToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{floatValue(7.11)})
	spec.Expect(filter("33", nil).(float64)).ToEqual(4.641350210970464)
}

func TestDivideByADynamicIntValue(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{dynamicValue("count")})
	spec.Expect(filter("33", params("count", 112)).(float64)).ToEqual(0.29464285714285715)
}

func TestDivideByADynamicFloatValue(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{dynamicValue("count")})
	spec.Expect(filter("12", params("count", 44.2)).(float64)).ToEqual(0.27149321266968324)
}

func TestDivideByDynamicNoop(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]core.Value{dynamicValue("count")})
	spec.Expect(filter("12", params("count", "22")).(string)).ToEqual("12")
}
