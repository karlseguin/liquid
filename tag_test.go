package liquid

import (
	"testing"
	"github.com/karlseguin/gspec"
)

func TestTagHandlesEmptyTag(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{}}"))
	spec.Expect(tag).ToBeNil()
	spec.Expect(err).ToBeNil()
}

func TestTagHandlesSpaceOnlyTag(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{   }}"))
	spec.Expect(tag).ToBeNil()
	spec.Expect(err).ToBeNil()
}

func TestTagExtractsASimpleStatic(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{  'over 9000'}}"))
	spec.Expect(string(tag.(*StaticTag).Value)).ToEqual("over 9000")
	spec.Expect(err).ToBeNil()
}

func TestTagExtractsAComplexStatic(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{'it\\'s over \\9000'}}"))
	spec.Expect(string(tag.(*StaticTag).Value)).ToEqual("it's over \\9000")
	spec.Expect(err).ToBeNil()
}

func TestTagExtractsAStaticWithAnEndingQuote(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{'it\\''}}"))
	spec.Expect(string(tag.(*StaticTag).Value)).ToEqual("it'")
	spec.Expect(err).ToBeNil()
}

func TestTagExtractionGivesErrorForUnclosedStatic(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{ 'failure }}"))
	spec.Expect(tag).ToBeNil()
	spec.Expect(err.Error()).ToEqual(`Tag had an unclosed single quote in "{{ 'failure }}"`)
}

func TestTagExtractASimpleDynamic(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{ name  }}"))
	assertDynamic(spec, tag, "name")
	spec.Expect(err).ToBeNil()
}

func TestTagExtractANestedDynamic(t *testing.T) {
	spec := gspec.New(t)
	tag, err := tagExtractor([]byte("{{ user.name.first  }}"))
	assertDynamic(spec, tag, "user", "name", "first")
	spec.Expect(err).ToBeNil()
}

func assertDynamic(spec *gspec.S, tag Token, expected ...string) {
	d := tag.(*DynamicTag)
	spec.Expect(len(d.Values)).ToEqual(len(expected))
	for index, e := range expected {
		spec.Expect(string(d.Values[index])).ToEqual(e)
	}
}
