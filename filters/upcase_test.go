package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestUpcasesAString(t *testing.T) {
	spec := gspec.New(t)
	filter := UpcaseFactory(nil)
	spec.Expect(filter("dbz", nil).(string)).ToEqual("DBZ")
}

func TestUpcasesBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := UpcaseFactory(nil)
	spec.Expect(string(filter([]byte("dbz"), nil).([]byte))).ToEqual("DBZ")
}

func TestUpcasesPassThroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := UpcaseFactory(nil)
	spec.Expect(filter(123, nil).(int)).ToEqual(123)
}
