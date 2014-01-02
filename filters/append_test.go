package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestAppendToAString(t *testing.T) {
	spec := gspec.New(t)
	filter := AppendFactory([]string{"?!"})
	spec.Expect(filter("dbz").(string)).ToEqual("dbz?!")
}

func TestAppendToBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := AppendFactory([]string{"boring"})
	spec.Expect(filter([]byte("so")).(string)).ToEqual("soboring")
}
