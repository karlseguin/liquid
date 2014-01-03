package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestDowncasesAString(t *testing.T) {
	spec := gspec.New(t)
	filter := DowncaseFactory(nil)
	spec.Expect(filter("DBZ", nil).(string)).ToEqual("dbz")
}

func TestDowncasesBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := DowncaseFactory(nil)
	spec.Expect(string(filter([]byte("DBZ"), nil).([]byte))).ToEqual("dbz")
}

func TestDowncasesPassThroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := DowncaseFactory(nil)
	spec.Expect(filter(123, nil).(int)).ToEqual(123)
}
