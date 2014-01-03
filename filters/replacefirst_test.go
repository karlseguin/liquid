package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestReplaceFirstValueInAString(t *testing.T) {
	spec := gspec.New(t)
	filter := ReplaceFirstFactory([]core.Value{stringValue("foo"), stringValue("bar")})
	spec.Expect(filter("foobarforfoo", nil).(string)).ToEqual("barbarforfoo")
}
