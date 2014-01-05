package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
	"time"
)

func TestPlusAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{intValue(5)})
	spec.Expect(filter(43, nil).(int)).ToEqual(48)
}

func TestPlusAFloattToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{floatValue(5.11)})
	spec.Expect(filter(43, nil).(float64)).ToEqual(48.11)
}

func TestPlusAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{intValue(5)})
	spec.Expect(filter(43.2, nil).(float64)).ToEqual(48.2)
}

func TestPlusAnIntToNow(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{intValue(61)})
	spec.Expect(filter("now", nil).(time.Time)).ToEqual(core.Now().Add(time.Minute * 61))
}

func TestPlusAnIntToATime(t *testing.T) {
	spec := gspec.New(t)
	now := time.Now()
	filter := PlusFactory([]core.Value{intValue(7)})
	spec.Expect(filter(now, nil).(time.Time)).ToEqual(now.Add(time.Minute * 7))
}

func TestPlusAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{intValue(7)})
	spec.Expect(filter("33", nil).(int)).ToEqual(40)
}

func TestPlusAnIntToAStringAsAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{floatValue(2.2)})
	spec.Expect(filter("33.11", nil).(float64)).ToEqual(35.31)
}

func TestPlusAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{intValue(7)})
	spec.Expect(filter([]byte("34"), nil).(int)).ToEqual(41)
}

func TestPlusAnDynamicIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{dynamicValue("fee")})
	spec.Expect(filter([]byte("34"), params("fee", 5)).(int)).ToEqual(39)
}

func TestPlusAnDynamicFloatToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]core.Value{dynamicValue("fee")})
	spec.Expect(filter([]byte("34"), params("fee", 5.1)).(float64)).ToEqual(39.1)
}
