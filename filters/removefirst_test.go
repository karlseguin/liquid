package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestRemovesFirstValueFromAString(t *testing.T) {
	spec := gspec.New(t)
	filter := RemoveFirstFactory([]core.Value{stringValue("foo")})
	spec.Expect(filter("foobarforfoo", nil).(string)).ToEqual("barforfoo")
}
