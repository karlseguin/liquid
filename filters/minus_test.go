package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
	"time"
)

func TestMinusAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"5"})
	spec.Expect(filter(43).(int)).ToEqual(38)
}

func TestMinusAnIntToAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"5"})
	spec.Expect(filter(43.2).(float64)).ToEqual(38.2)
}

func TestMinusAnIntToATime(t *testing.T) {
	spec := gspec.New(t)
	now := time.Now()
	filter := MinusFactory([]string{"7"})
	spec.Expect(filter(now).(time.Time)).ToEqual(now.Add(time.Minute * -7))
}

func TestMinusAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"7"})
	spec.Expect(filter("33").(int)).ToEqual(26)
}

func TestMinusAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"7"})
	spec.Expect(filter([]byte("34")).(int)).ToEqual(27)
}

func TestMinusAnIntToAStringAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"7"})
	spec.Expect(filter("abc").(string)).ToEqual("abc")
}

func TestMinusAnIntToBytesAsAString(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"8"})
	spec.Expect(filter([]byte("abb")).(string)).ToEqual("abb")
}

func TestMinusAStringToBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"xx"})
	spec.Expect(filter([]byte("abb")).(string)).ToEqual("abb")
}

func TestMinusAStringToString(t *testing.T) {
	spec := gspec.New(t)
	filter := MinusFactory([]string{"axx9"})
	spec.Expect(filter("dasdb").(string)).ToEqual("dasdb")
}
