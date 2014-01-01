package liquid

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestParsesATextOnlyTemplate(t *testing.T) {
	spec := gspec.New(t)
	template, _ := ParseString("it's over 9000", nil)
	spec.Expect(len(template.Tokens)).ToEqual(1)
	assertLiteral(t, template, 0, "it's over 9000")
}

func assertLiteral(t *testing.T, template *Template, index int, expected string) {
	actual := string(template.Tokens[index].(*Literal).Value)
	if actual != expected {
		t.Errorf("Expected token %d to be a literal with value %q, got %q", index, expected, actual)
	}
}
