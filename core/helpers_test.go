package core

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestSkipSpaceHandlesValueWithOnlySpaces(t *testing.T) {
	spec := gspec.New(t)
	spec.Expect(SkipSpaces([]byte("    "))).ToEqual(-1)
}

func TestSkipSpaceReturnsTheIndexOfTheFirstNonSpace(t *testing.T) {
	spec := gspec.New(t)
	spec.Expect(SkipSpaces([]byte("   over 9000"))).ToEqual(3)
}

func TestSkipSpaceWithNoSpaces(t *testing.T) {
	spec := gspec.New(t)
	spec.Expect(SkipSpaces([]byte("over 9000"))).ToEqual(0)
}

func TestTrimSringArrayWhenCapIsLen(t *testing.T) {
	actual := TrimStrings([]string{"it's", "over", "9000"})
	assertStringArray(t, actual, "it's", "over", "9000")
}

func TestTrimSringArrayWhenCapIsLarger(t *testing.T) {
	arr := make([]string, 0)
	arr = append(arr, "it's")
	arr = append(arr, "over")
	arr = append(arr, "9000")
	actual := TrimStrings(arr)
	assertStringArray(t, actual, "it's", "over", "9000")
}

func TestToBytesForString(t *testing.T) {
	assertBytes(t, ToBytes("it's over 9000"), "it's over 9000")
}

func TestToBytesForBytes(t *testing.T) {
	assertBytes(t, ToBytes([]byte("it's over 9000")), "it's over 9000")
}

func TestToBytesForInt(t *testing.T) {
	assertBytes(t, ToBytes(9000), "9000")
}

func TestToBytesForFloat(t *testing.T) {
	assertBytes(t, ToBytes(9000.132), "9000.132")
}

func TestToBytesForBool(t *testing.T) {
	assertBytes(t, ToBytes(true), "true")
	assertBytes(t, ToBytes(false), "false")
}

func TestToBytesForStringer(t *testing.T) {
	assertBytes(t, ToBytes(new(Stringable)), "i am a stringer")
}

func TestToBytesForANoneStringer(t *testing.T) {
	assertBytes(t, ToBytes(&NotStringable{"leto", 1400}), "&{leto 1400}")
}

func assertStringArray(t *testing.T, actuals []string, expected ...string) {
	spec := gspec.New(t)
	spec.Expect(len(actuals)).ToEqual(cap(actuals))
	spec.Expect(len(actuals)).ToEqual(len(expected))
	for i, a := range actuals {
		spec.Expect(a).ToEqual(expected[i])
	}
}

func assertBytes(t *testing.T, actual []byte, expected string) {
	if string(actual) != expected {
		t.Errorf("Expected %q to equal %q", string(actual), expected)
	}
}

type Stringable struct {
}

func (s *Stringable) String() string {
	return "i am a stringer"
}

type NotStringable struct {
	Name string
	Age  int
}
