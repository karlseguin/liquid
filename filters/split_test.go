package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestSplitsAStringOnDefaultSpace(t *testing.T) {
	spec := gspec.New(t)
	filter := SplitFactory([]core.Value{})
	values := filter("hello world", nil).([]string)
	spec.Expect(len(values)).ToEqual(2)
	spec.Expect(values[0]).ToEqual("hello")
	spec.Expect(values[1]).ToEqual("world")
}

func TestSplitsAStringOnSpecifiedValue(t *testing.T) {
	spec := gspec.New(t)
	filter := SplitFactory([]core.Value{stringValue("..")})
	values := filter([]byte("hel..lowo..rl..d"), nil).([]string)
	spec.Expect(len(values)).ToEqual(4)
	spec.Expect(values[0]).ToEqual("hel")
	spec.Expect(values[1]).ToEqual("lowo")
	spec.Expect(values[2]).ToEqual("rl")
	spec.Expect(values[3]).ToEqual("d")
}

func TestSplitsAStringOnADynamicValue(t *testing.T) {
	spec := gspec.New(t)
	filter := SplitFactory([]core.Value{dynamicValue("sep")})
	values := filter("over;9000;!", params("sep", ";")).([]string)
	spec.Expect(len(values)).ToEqual(3)
	spec.Expect(values[0]).ToEqual("over")
	spec.Expect(values[1]).ToEqual("9000")
	spec.Expect(values[2]).ToEqual("!")
}
