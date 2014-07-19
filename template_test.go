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

func TestRendersOutputTagsWithMap(t *testing.T) {
	d := map[string]interface{}{
		"ghola": map[string]interface{}{"incarnations": 67, "master": "LETO"},
	}
	template, _ := ParseString("duncan, next is {{ ghola.incarnations | plus: 1}}th. Your master is {{ ghola.master | upcase }}", nil)
	assertRender(t, template, d, `duncan, next is 68th. Your master is LETO`)
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

var complexTemplate = `
Out of
{% for color in colors reverse %}
- {{ color}}
{% endfor %}
{% capture favorite %}{{ colors |first}}{%endcapture%}
Your favorite color was {{favorite}}.
---
{% if ghola.incarnations > 10%}
You've been raised many times
{%else   %}
Youngn'
{%endif%}
---
{% for i in ( 1 ..ghola.name.size)%}
{%case i%}
{%when 2%}{%   continue%}
{%when 4%}{%   break%}
{%   endcase   %}
{{ i | minus:1 }} is {{ ghola.name[i]}}
{% endfor %}`

func TestTemplateRender1(t *testing.T) {
	d := map[string]interface{}{
		"ghola":  PersonS{"Duncan", 5},
		"colors": []string{"blue", "red", "white"},
	}
	template, _ := ParseString(complexTemplate, nil)
	assertRender(t, template, d, "\nOut of\n\n- white \n- red \n- blue  \nYour favorite color was blue.\n---\n\nYoungn'\n\n---\n  0 is 68    2 is 110  ")
}

func BenchmarkParseTemplateWithoutCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseString(complexTemplate, NoCache)
	}
}

func BenchmarkParseTemplateWithCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseString(complexTemplate, nil)
	}
}

func BenchmarkRenderTemplate(b *testing.B) {
	d := map[string]interface{}{
		"ghola":  PersonS{"Duncan", 5},
		"colors": []string{"blue", "red", "white"},
	}
	template, _ := ParseString(complexTemplate, nil)
	writer := new(NilWriter)
	for i := 0; i < b.N; i++ {
		template.Render(writer, d)
	}
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

type NilWriter struct{}

func (w *NilWriter) Write(b []byte) (int, error) {
	return 0, nil
}
