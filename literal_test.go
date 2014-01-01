package liquid

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestLiteralExtractorCreatesACopy(t *testing.T) {
	spec := gspec.New(t)
	original := []byte("it's over 9000")
	token, _ := literalExtractor(original)
	original[10] = '8'
	spec.Expect(string(token.(*Literal).Value)).ToEqual("it's over 9000")
}
