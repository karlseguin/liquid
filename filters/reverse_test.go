package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestReverseDoesNothingOnInvalidType(t *testing.T) {
	spec := gspec.New(t)
	filter := ReverseFactory(nil)
	spec.Expect(filter(123, nil).(int)).ToEqual(123)
}

func TestReverseAnEvenLengthString(t *testing.T) {
	spec := gspec.New(t)
	filter := ReverseFactory(nil)
	spec.Expect(string(filter("123456", nil).([]byte))).ToEqual("654321")
}

func TestReverseAnOddLengthString(t *testing.T) {
	spec := gspec.New(t)
	filter := ReverseFactory(nil)
	spec.Expect(string(filter("12345", nil).([]byte))).ToEqual("54321")
}

func TestReverseASingleCharacterString(t *testing.T) {
	spec := gspec.New(t)
	filter := ReverseFactory(nil)
	spec.Expect(string(filter("1", nil).([]byte))).ToEqual("1")
}

func TestReverseAnEvenLengthArray(t *testing.T) {
	spec := gspec.New(t)
	filter := ReverseFactory(nil)
	values := filter([]int{1, 2, 3, 4}, nil).([]int)
	spec.Expect(len(values)).ToEqual(4)
	spec.Expect(values[0]).ToEqual(4)
	spec.Expect(values[1]).ToEqual(3)
	spec.Expect(values[2]).ToEqual(2)
	spec.Expect(values[3]).ToEqual(1)
}

func TestReverseAnOddLengthArray(t *testing.T) {
	spec := gspec.New(t)
	filter := ReverseFactory(nil)
	values := filter([]float64{1.1, 2.2, 3.3}, nil).([]float64)
	spec.Expect(len(values)).ToEqual(3)
	spec.Expect(values[0]).ToEqual(3.3)
	spec.Expect(values[1]).ToEqual(2.2)
	spec.Expect(values[2]).ToEqual(1.1)
}

func TestReverseASingleElementArray(t *testing.T) {
	spec := gspec.New(t)
	filter := ReverseFactory(nil)
	values := filter([]bool{true}, nil).([]bool)
	spec.Expect(len(values)).ToEqual(1)
	spec.Expect(values[0]).ToEqual(true)
}
