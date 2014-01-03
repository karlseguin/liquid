package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestJoinsStringsWithTheSpecifiedGlue(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory([]core.Value{stringValue("..")})
	spec.Expect(string(filter([]string{"leto", "atreides"}, nil).([]byte))).ToEqual("leto..atreides")
}

func TestJoinsVariousTypesWithTheDefaultGlue(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	spec.Expect(string(filter([]interface{}{"leto", 123, true}, nil).([]byte))).ToEqual("leto 123 true")
}

func TestJoinPassthroughOnEmptyArray(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	arr := [0]int{}
	spec.Expect(filter(arr, nil).([0]int)).ToEqual(arr)
}

func TestJoinPassthroughOnEmptySlice(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	arr := []int{}
	spec.Expect(len(filter(arr, nil).([]int))).ToEqual(0)
}

func TestJoinPassthroughOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := JoinFactory(nil)
	spec.Expect(filter("hahah", nil).(string)).ToEqual("hahah")
}
