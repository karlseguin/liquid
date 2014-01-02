package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestReplaceFirstValueInAString(t *testing.T) {
	spec := gspec.New(t)
	filter := ReplaceFirstFactory([]string{"foo", "bar"})
	spec.Expect(filter("foobarforfoo").(string)).ToEqual("barbarforfoo")
}
