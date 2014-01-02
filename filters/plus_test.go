package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
	"time"
)

func TestPlusAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"5"})
	spec.Expect(filter(43).(int)).ToEqual(48)
}

func TestPlusAFloattToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"5.11"})
	spec.Expect(filter(43).(float64)).ToEqual(48.11)
}

func TestPlusAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"5"})
	spec.Expect(filter(43.2).(float64)).ToEqual(48.2)
}

func TestPlusAnIntToATime(t *testing.T) {
	spec := gspec.New(t)
	now := time.Now()
	filter := PlusFactory([]string{"7"})
	spec.Expect(filter(now).(time.Time)).ToEqual(now.Add(time.Minute * 7))
}

func TestPlusAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"7"})
	spec.Expect(filter("33").(int)).ToEqual(40)
}

func TestPlusAnIntToAStringAsAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"2.2"})
	spec.Expect(filter("33.11").(float64)).ToEqual(35.31)
}

func TestPlusAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"7"})
	spec.Expect(filter([]byte("34")).(int)).ToEqual(41)
}

func TestPlusAnIntToAStringAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"7"})
	spec.Expect(filter("abc").(string)).ToEqual("abc7")
}

func TestPlusAnIntToBytesAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"8"})
	spec.Expect(filter([]byte("abb")).(string)).ToEqual("abb8")
}

func TestPlusAStringToBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"xx"})
	spec.Expect(filter([]byte("abb")).(string)).ToEqual("abbxx")
}

func TestPlusAStringToString(t *testing.T) {
	spec := gspec.New(t)
	filter := PlusFactory([]string{"axx9"})
	spec.Expect(filter("dasdb").(string)).ToEqual("dasdbaxx9")
}
