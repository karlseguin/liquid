package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestRemovesFirstValueFromAString(t *testing.T) {
	spec := gspec.New(t)
	filter := RemoveFirstFactory([]string{"foo"})
	spec.Expect(filter("foobarforfoo").(string)).ToEqual("barforfoo")
}
