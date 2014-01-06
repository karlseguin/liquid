package liquid

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestParsesATextOnlyTemplate(t *testing.T) {
	spec := gspec.New(t)
	template, _ := ParseString("it's over 9000", nil)
	spec.Expect(len(template.Code)).ToEqual(1)
	assertLiteral(t, template, 0, "it's over 9000")
}

func TestRendersOutputTags(t *testing.T) {
	d := map[string]interface{}{
		"name":   "leto atreides",
		"colors": []string{"brown", "blue"},
	}
	template, _ := ParseString("hello {{name | capitalize }}, you ranked {{ colors | first }} as your favorite color", nil)
	assertRender(t, template, d, `hello Leto Atreides, you ranked brown as your favorite color`)
}

func TestRendersOutputTagsWithStructPointers(t *testing.T) {
	d := map[string]interface{}{
		"ghola": &Person{"Duncan", 67, &Person{"Leto", 0, nil}},
	}
	template, _ := ParseString("{{ ghola | downcase }}, next is {{ ghola.incarnations | plus: 1}}th. Your master is {{ ghola.master | upcase }}", nil)
	assertRender(t, template, d, `duncan, next is 68th. Your master is LETO`)
}

func TestRendersOutputTagsWithStructs(t *testing.T) {
	d := map[string]interface{}{
		"ghola": PersonS{"Duncan", 67},
	}
	template, _ := ParseString("{{ ghola | downcase }}, next is {{ ghola.incarnations | plus: 1}}th. Your master is {{ ghola.master | upcase }}", nil)
	assertRender(t, template, d, `duncan, next is 68th. Your master is {{GHOLA.MASTER}}`)
}

func TestRendersCaptureOfSimpleText(t *testing.T) {
	d := map[string]interface{}{
		"ghola": PersonS{"Duncan", 67},
	}
	template, _ := ParseString("welcome {% capture intro %}Mr.X{%  endcapture%}. {{ intro }}", nil)
	assertRender(t, template, d, `welcome . Mr.X`)
}

func TestRendersCaptureWithNestedOutputs(t *testing.T) {
	d := map[string]interface{}{
		"ghola": PersonS{"Duncan", 67},
	}
	template, _ := ParseString("welcome{%   capture name   %} {{ ghola | downcase }}{%endcapture%}! {{ name }}", nil)
	assertRender(t, template, d, `welcome!  duncan`)
}

func TestRenderSimpleIfstatement(t *testing.T) {
	template, _ := ParseString("A-{% if 2 == 2 %}in if{% endif %}-Z", nil)
	assertRender(t, template, nil, `A-in if-Z`)
}

func TestRenderSimpleElseIfstatement(t *testing.T) {
	template, _ := ParseString("A-{% if 0 == 2 %}in if{% elseif 2 == 2 %}in elseif{% endif %}-Z", nil)
	assertRender(t, template, nil, `A-in elseif-Z`)
}

func TestRenderSimpleElseStatement(t *testing.T) {
	template, _ := ParseString("A-{% if 0 == 2 %}in if{% elseif 2 == 0 %}in elseif{% else %}in else{% endif %}-Z", nil)
	assertRender(t, template, nil, `A-in else-Z`)
}

func TestRenderANilCheckAgainstDynamicValue(t *testing.T) {
	d := map[string]interface{}{
		"ghola": PersonS{"Duncan", 67},
	}
	template, _ := ParseString("A-{% if false %}in if{% elseif ghola %}in elseif{% else %}in else{% endif %}-Z", nil)
	assertRender(t, template, d, `A-in elseif-Z`)
}

func TestRendersNothingForAFailedUnless(t *testing.T) {
	template, _ := ParseString("A-{% unless true %}in unless{%endunless%}-Z", nil)
	assertRender(t, template, nil, `A--Z`)
}

func TestRendersAnUnlessTag(t *testing.T) {
	template, _ := ParseString("A-{% unless false %}in unless{%endunless%}-Z", nil)
	assertRender(t, template, nil, `A-in unless-Z`)
}

func TestRendersElseAFailedUnless(t *testing.T) {
	template, _ := ParseString("A-{% unless true %}in if{%else%}in else{%endunless%}-Z", nil)
	assertRender(t, template, nil, `A-in else-Z`)
}

func TestRendersCaseWhen1(t *testing.T) {
	template, _ := ParseString("A-{% case 'abc' %}{% when 'abc' %}when1{% when 1 or 123 %}when2{% else %}else{% endcase%}-Z", nil)
	assertRender(t, template, nil, `A-when1-Z`)
}

func TestRendersCaseWhen2(t *testing.T) {
	template, _ := ParseString("A-{% case 123 %}{% when 'abc' %}when1{% when 1 or 123 %}when2{% else %}else{% endcase%}-Z", nil)
	assertRender(t, template, nil, `A-when2-Z`)
}

func TestRendersCaseElse(t *testing.T) {
	template, _ := ParseString("A-{% case other %}{% when 'abc' %}when1{% when 1 or 123 %}when2{% else %}else{% endcase%}-Z", nil)
	assertRender(t, template, nil, `A-else-Z`)
}

func assertLiteral(t *testing.T, template *Template, index int, expected string) {
	actual := string(template.Code[index].(*Literal).Value)
	if actual != expected {
		t.Errorf("Expected code %d to be a literal with value %q, got %q", index, expected, actual)
	}
}

type Person struct {
	Name         string
	Incarnations int
	Master       *Person
}

func (p *Person) String() string {
	return p.Name
}

type PersonS struct {
	Name         string
	Incarnations int
}

func (p PersonS) String() string {
	return p.Name
}
