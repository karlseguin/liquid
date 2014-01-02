package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestReplaceValuesInAString(t *testing.T) {
	spec := gspec.New(t)
	filter := ReplaceFactory([]string{"foo", "bar"})
	spec.Expect(filter("foobarforfoo").(string)).ToEqual("barbarforbar")
}
