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
	spec := gspec.New(t)
	d := map[string]interface{}{
		"name":   "leto atreides",
		"colors": []string{"brown", "blue"},
	}
	template, _ := ParseString("hello {{ name | capitalize }}, you ranked {{ colors | first }} as your favorite color", nil)
	spec.Expect(string(template.Render(d))).ToEqual(`hello Leto Atreides, you ranked brown as your favorite color`)
}

func TestRendersOutputTagsWithStructPointers(t *testing.T) {
	spec := gspec.New(t)
	d := map[string]interface{}{
		"ghola": &Person{"Duncan", 67, &Person{"Leto", 0, nil}},
	}
	template, _ := ParseString("{{ ghola | downcase }}, next is {{ ghola.incarnations | plus: 1}}th. Your master is {{ ghola.master | upcase }}", nil)
	spec.Expect(string(template.Render(d))).ToEqual(`duncan, next is 68th. Your master is LETO`)
}

func TestRendersOutputTagsWithStructs(t *testing.T) {
	spec := gspec.New(t)
	d := map[string]interface{}{
		"ghola": PersonS{"Duncan", 67},
	}
	template, _ := ParseString("{{ ghola | downcase }}, next is {{ ghola.incarnations | plus: 1}}th. Your master is {{ ghola.master | upcase }}", nil)
	spec.Expect(string(template.Render(d))).ToEqual(`duncan, next is 68th. Your master is {{GHOLA.MASTER}}`)
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
