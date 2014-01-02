package filters

import (
	"testing"
	"github.com/karlseguin/gspec"
)

func TestDowncasesAString(t *testing.T) {
	spec := gspec.New(t)
	filter := DowncaseFactory(nil)
	spec.Expect(filter("DBZ").(string)).ToEqual("dbz")
}

func TestDowncasesBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := DowncaseFactory(nil)
	spec.Expect(string(filter([]byte("DBZ")).([]byte))).ToEqual("dbz")
}

func TestDowncasesPassThroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := DowncaseFactory(nil)
	spec.Expect(filter(123).(int)).ToEqual(123)
}
