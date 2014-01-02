package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestRemovesValuesFromAString(t *testing.T) {
	spec := gspec.New(t)
	filter := RemoveFactory([]string{"foo"})
	spec.Expect(filter("foobarforfoo").(string)).ToEqual("barfor")
}
