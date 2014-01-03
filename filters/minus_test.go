package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
	"time"
)

func TestMinusAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{intValue(5)})
	spec.Expect(filter(43, nil).(int)).ToEqual(38)
}

func TestMinusAFloattToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{floatValue(5.11)})
	spec.Expect(filter(43, nil).(float64)).ToEqual(37.89)
}

func TestMinusAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{intValue(5)})
	spec.Expect(filter(43.2, nil).(float64)).ToEqual(38.2)
}

func TestMinusAnIntToATime(t *testing.T) {
	spec := gspec.New(t)
	now := time.Now()
	filter := MinusFactory([]core.Value{intValue(7)})
	spec.Expect(filter(now, nil).(time.Time)).ToEqual(now.Add(time.Minute * -7))
}

func TestMinusAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{intValue(7)})
	spec.Expect(filter("33", nil).(int)).ToEqual(26)
}

func TestMinusAnIntToAStringAsAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{floatValue(2.2)})
	spec.Expect(filter("33.11", nil).(float64)).ToEqual(30.91)
}

func TestMinusAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{intValue(7)})
	spec.Expect(filter([]byte("34"), nil).(int)).ToEqual(27)
}

func TestMinusAnDynamicIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{dynamicValue("fee")})
	spec.Expect(filter([]byte("34"), params("fee", 5)).(int)).ToEqual(29)
}

func TestMinusAnDynamicFloatToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]core.Value{dynamicValue("fee")})
	spec.Expect(filter([]byte("34"), params("fee", 5.1)).(float64)).ToEqual(28.9)
}
