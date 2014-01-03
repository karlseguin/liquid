package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestSizeOfString(t *testing.T) {
	spec := gspec.New(t)
	filter := SizeFactory(nil)
	spec.Expect(filter("dbz", nil).(int)).ToEqual(3)
}

func TestSizeOfByteArray(t *testing.T) {
	spec := gspec.New(t)
	filter := SizeFactory(nil)
	spec.Expect(filter([]byte("7 123"), nil).(int)).ToEqual(5)
}

func TestSizeOfIntArray(t *testing.T) {
	spec := gspec.New(t)
	filter := SizeFactory(nil)
	spec.Expect(filter([]int{2, 4, 5, 6}, nil).(int)).ToEqual(4)
}

func TestSizeOfBoolArray(t *testing.T) {
	spec := gspec.New(t)
	filter := SizeFactory(nil)
	spec.Expect(filter([]bool{true, false, true, true, false}, nil).(int)).ToEqual(5)
}

func TestSizeOfMap(t *testing.T) {
	spec := gspec.New(t)
	filter := SizeFactory(nil)
	spec.Expect(filter(map[string]int{"over": 9000}, nil).(int)).ToEqual(1)
}

func TestSizeOfSometingInvalid(t *testing.T) {
	spec := gspec.New(t)
	filter := SizeFactory(nil)
	spec.Expect(filter(false, nil).(bool)).ToEqual(false)
}
