package filters

import (
	"testing"
	"github.com/karlseguin/gspec"
)

func TestUpcasesAString(t *testing.T) {
	spec := gspec.New(t)
	filter := UpcaseFactory(nil)
	spec.Expect(filter("dbz").(string)).ToEqual("DBZ")
}

func TestUpcasesBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := UpcaseFactory(nil)
	spec.Expect(string(filter([]byte("dbz")).([]byte))).ToEqual("DBZ")
}

func TestUpcasesPassThroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := UpcaseFactory(nil)
	spec.Expect(filter(123).(int)).ToEqual(123)
}
