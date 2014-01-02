package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestReturnsTheLastItem(t *testing.T) {
	spec := gspec.New(t)
	filter := LastFactory(nil)
	spec.Expect(filter([]string{"leto", "atreides"}).(string)).ToEqual("atreides")
}

func TestReturnsTheLastItemIfOnlyOneItem(t *testing.T) {
	spec := gspec.New(t)
	filter := LastFactory(nil)
	spec.Expect(filter([]string{"leto"}).(string)).ToEqual("leto")
}

func TestReturnsTheLastItemOfAnArray(t *testing.T) {
	spec := gspec.New(t)
	filter := LastFactory(nil)
	arr := [4]int{1, 2, 3, 48}
	spec.Expect(filter(arr).(int)).ToEqual(48)
}

func TestLastPassthroughOnEmptyArray(t *testing.T) {
	spec := gspec.New(t)
	filter := LastFactory(nil)
	arr := [0]int{}
	spec.Expect(filter(arr).([0]int)).ToEqual(arr)
}

func TestLastPassthroughOnEmptySlice(t *testing.T) {
	spec := gspec.New(t)
	filter := LastFactory(nil)
	arr := []int{}
	spec.Expect(len(filter(arr).([]int))).ToEqual(0)
}

func TestLastPassthroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := LastFactory(nil)
	spec.Expect(filter("hahah").(string)).ToEqual("hahah")
}
