package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestTimesAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"5"})
	spec.Expect(filter(43).(int)).ToEqual(215)
}

func TestTimesAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"5"})
	spec.Expect(filter(43.3).(float64)).ToEqual(216.5)
}

func TestTimesAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"7"})
	spec.Expect(filter("33").(int)).ToEqual(231)
}

func TestTimesAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"7"})
	spec.Expect(filter([]byte("34")).(int)).ToEqual(238)
}

func TestTimesAnIntToAStringAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"7"})
	spec.Expect(filter("abc").(string)).ToEqual("abc")
}

func TestTimesAnIntToBytesAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"8"})
	spec.Expect(filter([]byte("abb")).(string)).ToEqual("abb")
}

func TestTimesAFloatToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"1.10"})
	spec.Expect(filter(43).(float64)).ToEqual(47.300000000000004)
}

func TestTimesAFloatToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"5.3"})
	spec.Expect(filter(43.3).(float64)).ToEqual(229.48999999999998)
}

func TestTimesAFloatToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := TimesFactory([]string{"7.11"})
	spec.Expect(filter("33").(float64)).ToEqual(234.63000000000002)
}
