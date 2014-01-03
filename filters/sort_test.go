package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestSortsAnArrayOfInteger(t *testing.T) {
	spec := gspec.New(t)
	filter := SortFactory(nil)
	values := filter([]int{3, 4, 1, 2, 3}, nil).([]int)
	spec.Expect(len(values)).ToEqual(5)
	spec.Expect(values[0]).ToEqual(1)
	spec.Expect(values[1]).ToEqual(2)
	spec.Expect(values[2]).ToEqual(3)
	spec.Expect(values[3]).ToEqual(3)
	spec.Expect(values[4]).ToEqual(4)
}

func TestSortsAnArrayOfStrings(t *testing.T) {
	spec := gspec.New(t)
	filter := SortFactory(nil)
	values := filter([]string{"cc", "b", "aa", "g"}, nil).([]string)
	spec.Expect(len(values)).ToEqual(4)
	spec.Expect(values[0]).ToEqual("aa")
	spec.Expect(values[1]).ToEqual("b")
	spec.Expect(values[2]).ToEqual("cc")
	spec.Expect(values[3]).ToEqual("g")
}

func TestSortsAnArrayOfFloats(t *testing.T) {
	spec := gspec.New(t)
	filter := SortFactory(nil)
	values := filter([]float64{1.1, 0.9, 1233.2, 21.994}, nil).([]float64)
	spec.Expect(len(values)).ToEqual(4)
	spec.Expect(values[0]).ToEqual(0.9)
	spec.Expect(values[1]).ToEqual(1.1)
	spec.Expect(values[2]).ToEqual(21.994)
	spec.Expect(values[3]).ToEqual(1233.2)
}

func TestSortsSortableData(t *testing.T) {
	spec := gspec.New(t)
	filter := SortFactory(nil)
	values := filter(People{&Person{"Leto"}, &Person{"Paul"}, &Person{"Jessica"}}, nil).(People)
	spec.Expect(len(values)).ToEqual(3)
	spec.Expect(values[0].Name).ToEqual("Jessica")
	spec.Expect(values[1].Name).ToEqual("Leto")
	spec.Expect(values[2].Name).ToEqual("Paul")
}

func TestSortsOtherValuesAsStrings(t *testing.T) {
	spec := gspec.New(t)
	filter := SortFactory(nil)
	values := filter([]interface{}{933, "spice", true, 123.44, "123", false}, nil).([]interface{})
	spec.Expect(len(values)).ToEqual(6)
	spec.Expect(values[0].(string)).ToEqual("123")
	spec.Expect(values[1].(float64)).ToEqual(123.44)
	spec.Expect(values[2].(int)).ToEqual(933)
	spec.Expect(values[3].(bool)).ToEqual(false)
	spec.Expect(values[4].(string)).ToEqual("spice")
	spec.Expect(values[5].(bool)).ToEqual(true)
}

func TestSortSkipsNonArrays(t *testing.T) {
	spec := gspec.New(t)
	filter := SortFactory(nil)
	spec.Expect(filter(1343, nil).(int)).ToEqual(1343)
}

type People []*Person

func (p People) Len() int {
	return len(p)
}

func (p People) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

func (p People) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Person struct {
	Name string
}

func (p *Person) String() string {
	return p.Name
}
