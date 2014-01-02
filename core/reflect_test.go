package core

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestResolvesAnInvalidValueFromNil(t *testing.T) {
	spec := gspec.New(t)
	spec.Expect(Resolve(nil, "something")).ToBeNil()
}

func TestResolvesAnInvalidValueFromAMap(t *testing.T) {
	spec := gspec.New(t)
	m := map[string]string{"name": "leto"}
	spec.Expect(Resolve(m, "age")).ToEqual("")
}

func TestResolvesAnInvalidValueFromAStruct(t *testing.T) {
	spec := gspec.New(t)
	m := &Person{"Leto", 3231}
	spec.Expect(Resolve(m, "IsGholas")).ToBeNil()
}

func TestResolvesAValueFromAMap(t *testing.T) {
	spec := gspec.New(t)
	m := map[string]string{"name": "leto"}
	spec.Expect(Resolve(m, "name")).ToEqual("leto")
}

func TestResolvesAValueFromAStruct(t *testing.T) {
	spec := gspec.New(t)
	m := &Person{"Leto", 3231}
	spec.Expect(Resolve(m, "name")).ToEqual("Leto")
	spec.Expect(Resolve(m, "age")).ToEqual(3231)
}

type Person struct {
	Name string
	Age  int
}
