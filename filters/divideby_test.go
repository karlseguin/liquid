package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestDivideByAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]string{"5"})
	spec.Expect(filter(43).(float64)).ToEqual(8.6)
}

func TestDivideByAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]string{"5.7"})
	spec.Expect(filter(43.3).(float64)).ToEqual(7.596491228070175)
}

func TestDivideByAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := DivideByFactory([]string{"7"})
	spec.Expect(filter("33").(float64)).ToEqual(4.714285714285714)
}
