package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestRemovesValuesFromAString(t *testing.T) {
	spec := gspec.New(t)
	filter := RemoveFactory([]core.Value{stringValue("foo")})
	spec.Expect(filter("foobarforfoo", nil).(string)).ToEqual("barfor")
}
