package liquid

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"strconv"
	"testing"
)

func TestOutputHandlesEmptyOutput(t *testing.T) {
	spec := gspec.New(t)
	output, err := newOutput(core.NewParser([]byte("{{}}")))
	spec.Expect(output).ToBeNil()
	spec.Expect(err).ToBeNil()
}

func TestOutputHandlesSpaceOnlyOutput(t *testing.T) {
	spec := gspec.New(t)
	output, err := newOutput(core.NewParser([]byte("{{   }}")))
	spec.Expect(output).ToBeNil()
	spec.Expect(err).ToBeNil()
}

func TestOutputExtractsASimpleStatic(t *testing.T) {
	output, _ := newOutput(core.NewParser([]byte("{{  'over 9000'}}")))
	assertRender(t, output, nil, "over 9000")
}

func TestOutputExtractsAComplexStatic(t *testing.T) {
	output, _ := newOutput(core.NewParser([]byte("{{'it\\'s over \\9000'}}")))
	assertRender(t, output, nil, "it's over \\9000")
}

func TestOutputExtractsAStaticWithAnEndingQuote(t *testing.T) {
	output, _ := newOutput(core.NewParser([]byte("{{'it\\''}}")))
	assertRender(t, output, nil, "it'")
}

func TestOutputExtractionGivesErrorForUnclosedStatic(t *testing.T) {
	spec := gspec.New(t)
	output, err := newOutput(core.NewParser([]byte("{{ 'failure }}")))
	spec.Expect(output).ToBeNil()
	spec.Expect(err.Error()).ToEqual(`Invalid value, a single quote might be missing ("{{ 'failure }}" - line 1)`)
}

func TestOutputNoFiltersForStatic(t *testing.T) {
	spec := gspec.New(t)
	output, _ := newOutput(core.NewParser([]byte("{{'fun'}}")))
	spec.Expect(len(output.(*Output).Filters)).ToEqual(0)
}

func TestOutputGeneratesErrorOnUnknownFilter(t *testing.T) {
	spec := gspec.New(t)
	_, err := newOutput(core.NewParser([]byte("{{'fun' | unknown }}")))
	spec.Expect(err.Error()).ToEqual(`Unknown filter "unknown" ("{{'fun' | unknown }}" - line 1)`)
}

func TestOutputGeneratesErrorOnInvalidParameter(t *testing.T) {
	spec := gspec.New(t)
	_, err := newOutput(core.NewParser([]byte("{{'fun' | debug: 'missing }}")))
	spec.Expect(err.Error()).ToEqual(`Invalid value, a single quote might be missing ("{{'fun' | debug: 'missing }}" - line 1)`)
}

func TestOutputWithASingleFilter(t *testing.T) {
	output, _ := newOutput(core.NewParser([]byte("{{'fun' | debug }}")))
	assertFilters(t, output, "debug(0)")
}

func TestOutputWithMultipleFilters(t *testing.T) {
	output, _ := newOutput(core.NewParser([]byte("{{'fun' | debug | debug}}")))
	assertFilters(t, output, "debug(0)", "debug(1)")
}

func TestOutputWithMultipleFiltersHavingParameters(t *testing.T) {
	spec := gspec.New(t)
	output, err := newOutput(core.NewParser([]byte("{{'fun' | debug:1,2 | debug:'test' | debug : 'test' , 5}}")))
	spec.Expect(err).ToBeNil()
	assertFilters(t, output, "debug(0, 1, 2)", "debug(1, test)", "debug(2, test, 5)")
}

func TestOutputWithAnEscapeParameter(t *testing.T) {
	spec := gspec.New(t)
	output, err := newOutput(core.NewParser([]byte("{{'fun' | debug: 'te\\'st'}}")))
	spec.Expect(err).ToBeNil()
	assertFilters(t, output, "debug(0, te'st)")
}

func assertFilters(t *testing.T, output core.Code, expected ...string) {
	spec := gspec.New(t)
	filters := output.(*Output).Filters
	spec.Expect(len(filters)).ToEqual(len(expected))
	for index, filter := range filters {
		actual := string(filter(strconv.Itoa(index), nil).([]byte))
		spec.Expect(actual).ToEqual(expected[index])
	}
}
