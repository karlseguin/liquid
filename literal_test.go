package liquid

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestLiteralRendersItself(t *testing.T) {
	spec := gspec.New(t)
	literal := newLiteral([]byte("it's over 9001"))
	spec.Expect(string(literal.Render(nil))).ToEqual("it's over 9001")
}
