package liquid

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestOutputHandlesEmptyOutput(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{}}"))
	spec.Expect(output).ToBeNil()
	spec.Expect(err).ToBeNil()
}

func TestOutputHandlesSpaceOnlyOutput(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{   }}"))
	spec.Expect(output).ToBeNil()
	spec.Expect(err).ToBeNil()
}

func TestOutputExtractsASimpleStatic(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{  'over 9000'}}"))
	spec.Expect(string(output.(*StaticOutput).Value)).ToEqual("over 9000")
	spec.Expect(err).ToBeNil()
}

func TestOutputExtractsAComplexStatic(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{'it\\'s over \\9000'}}"))
	spec.Expect(string(output.(*StaticOutput).Value)).ToEqual("it's over \\9000")
	spec.Expect(err).ToBeNil()
}

func TestOutputExtractsAStaticWithAnEndingQuote(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{'it\\''}}"))
	spec.Expect(string(output.(*StaticOutput).Value)).ToEqual("it'")
	spec.Expect(err).ToBeNil()
}

func TestOutputExtractionGivesErrorForUnclosedStatic(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{ 'failure }}"))
	spec.Expect(output).ToBeNil()
	spec.Expect(err.Error()).ToEqual(`Output had an unclosed single quote in "{{ 'failure }}"`)
}

func TestOutputExtractASimpleDynamic(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{ name  }}"))
	assertDynamic(spec, output, "name")
	spec.Expect(err).ToBeNil()
}

func TestOutputExtractANestedDynamic(t *testing.T) {
	spec := gspec.New(t)
	output, err := outputExtractor([]byte("{{ user.name.first  }}"))
	assertDynamic(spec, output, "user", "name", "first")
	spec.Expect(err).ToBeNil()
}

func assertDynamic(spec *gspec.S, output Token, expected ...string) {
	d := output.(*DynamicOutput)
	spec.Expect(len(d.Values)).ToEqual(len(expected))
	for index, e := range expected {
		spec.Expect(string(d.Values[index])).ToEqual(e)
	}
}
