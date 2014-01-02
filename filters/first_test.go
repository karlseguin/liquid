package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestReturnsTheFirstItem(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	spec.Expect(filter([]string{"leto", "atreides"}).(string)).ToEqual("leto")
}

func TestReturnsTheFirstItemIfOnlyOneItem(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	spec.Expect(filter([]string{"leto"}).(string)).ToEqual("leto")
}

func TestReturnsTheFirstItemOfAnArray(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	arr := [4]int{12, 2, 3, 48}
	spec.Expect(filter(arr).(int)).ToEqual(12)
}

func TestFirstPassthroughOnEmptyArray(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	arr := [0]int{}
	spec.Expect(filter(arr).([0]int)).ToEqual(arr)
}

func TestFirstPassthroughOnEmptySlice(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	arr := []int{}
	spec.Expect(len(filter(arr).([]int))).ToEqual(0)
}

func TestFirstPassthroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	spec.Expect(filter("hahah").(string)).ToEqual("hahah")
}
