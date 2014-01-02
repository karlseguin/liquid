package filters

import (
	"testing"
	"github.com/karlseguin/gspec"
)

func TestJoinsStringsWithTheSpecifiedGlue(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory([]string{".."})
	spec.Expect(string(filter([]string{"leto", "atreides"}).([]byte))).ToEqual("leto..atreides")
}

func TestJoinsVariousTypesWithTheDefaultGlue(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	spec.Expect(string(filter([]interface{}{"leto", 123, true}).([]byte))).ToEqual("leto 123 true")
}

func TestJoinPassthroughOnEmptyArray(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	arr := [0]int{}
	spec.Expect(filter(arr).([0]int)).ToEqual(arr)
}

func TestJoinPassthroughOnEmptySlice(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	arr := []int{}
	spec.Expect(len(filter(arr).([]int))).ToEqual(0)
}

func TestJoinPassthroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	spec.Expect(filter("hahah").(string)).ToEqual("hahah")
}
