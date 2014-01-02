package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestCapitalizesAString(t *testing.T) {
	spec := gspec.New(t)
	filter := CapitalizeFactory(nil)
	spec.Expect(string(filter("tiger got to hunt, bird got to fly").([]byte))).ToEqual("Tiger Got To Hunt, Bird Got To Fly")
}

func TestCapitalizesBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := CapitalizeFactory(nil)
	spec.Expect(string(filter([]byte("Science is magic that works ")).([]byte))).ToEqual("Science Is Magic That Works ")
}

func TestCapitalizePassThroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := DowncaseFactory(nil)
	spec.Expect(filter(123).(int)).ToEqual(123)
}
