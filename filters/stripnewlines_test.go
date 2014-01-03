package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestStripsNewLinesFromStirng(t *testing.T) {
	spec := gspec.New(t)
	filter := StripNewLinesFactory(nil)
	spec.Expect(filter("f\no\ro\n\r", nil).(string)).ToEqual("foo")
}
