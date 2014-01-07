package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"strings"
	"testing"
)

func TestTimesAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{intValue(5)})
	spec.Expect(filter(43, nil).(int)).ToEqual(215)
}

func TestTimesAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{intValue(2)})
	spec.Expect(filter(43.3, nil).(float64)).ToEqual(86.6)
}

func TestTimesAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{intValue(7)})
	spec.Expect(filter("33", nil).(int)).ToEqual(231)
}

func TestTimesAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{intValue(7)})
	spec.Expect(filter([]byte("34"), nil).(int)).ToEqual(238)
}

func TestTimesAnIntToAStringAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{intValue(7)})
	spec.Expect(filter("abc", nil).(string)).ToEqual("abc")
}

func TestTimesAnIntToBytesAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{intValue(8)})
	spec.Expect(filter([]byte("abb"), nil).(string)).ToEqual("abb")
}

func TestTimesAFloatToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{floatValue(1.10)})
	spec.Expect(filter(43, nil).(float64)).ToEqual(47.300000000000004)
}

func TestTimesAFloatToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{floatValue(5.3)})
	spec.Expect(filter(43.3, nil).(float64)).ToEqual(229.48999999999998)
}

func TestTimesAFloatToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{floatValue(7.11)})
	spec.Expect(filter("33", nil).(float64)).ToEqual(234.63000000000002)
}

func TestTimesADynamicIntValue(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{dynamicValue("count")})
	spec.Expect(filter("33", params("count", 112)).(int)).ToEqual(3696)
}

func TestTimesADynamicFloatValue(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{dynamicValue("count")})
	spec.Expect(filter("12", params("count", 44.2)).(float64)).ToEqual(530.4000000000001)
}

func TestTimesDynamicNoop(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]core.Value{dynamicValue("count")})
	spec.Expect(filter("12", params("count", "22")).(string)).ToEqual("12")
}

func stringValue(s string) core.Value {
	return &core.StaticStringValue{s}
}

func intValue(n int) core.Value {
	return &core.StaticIntValue{n}
}

func floatValue(f float64) core.Value {
	return &core.StaticFloatValue{f}
}

func dynamicValue(s string) core.Value {
	return core.NewDynamicValue(strings.Split(s, "."))
}

func params(values ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1]
	}
	return m
}
