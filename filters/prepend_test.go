package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestPrependToAString(t *testing.T) {
	spec := gspec.New(t)
	filter := PrependFactory([]string{"?!"})
	spec.Expect(filter("dbz").(string)).ToEqual("?!dbz")
}

func TestPrependToBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := PrependFactory([]string{"boring"})
	spec.Expect(filter([]byte("so")).(string)).ToEqual("boringso")
}
