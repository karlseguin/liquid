package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestReturnsTheFirstItem(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	spec.Expect(filter([]string{"leto", "atreides"}, nil).(string)).ToEqual("leto")
}

func TestReturnsTheFirstItemIfOnlyOneItem(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	spec.Expect(filter([]string{"leto"}, nil).(string)).ToEqual("leto")
}

func TestReturnsTheFirstItemOfAnArray(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	arr := [4]int{12, 2, 3, 48}
	spec.Expect(filter(arr, nil).(int)).ToEqual(12)
}

func TestFirstPassthroughOnEmptyArray(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	arr := [0]int{}
	spec.Expect(filter(arr, nil).([0]int)).ToEqual(arr)
}

func TestFirstPassthroughOnEmptySlice(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	arr := []int{}
	spec.Expect(len(filter(arr, nil).([]int))).ToEqual(0)
}

func TestFirstPassthroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := FirstFactory(nil)
	spec.Expect(filter("hahah", nil).(string)).ToEqual("hahah")
}
