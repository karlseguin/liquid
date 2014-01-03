package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestTruncateAString(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{intValue(3)})
	spec.Expect(filter("1234567", nil).(string)).ToEqual("123")
}

func TestTruncateAShortString(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{intValue(100)})
	spec.Expect(filter("1234567", nil).(string)).ToEqual("1234567")
}

func TestTruncateAPerfectString(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{intValue(7)})
	spec.Expect(filter("1234567", nil).(string)).ToEqual("1234567")
}

func TestTruncateAnAlmostPerfectString(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{intValue(6)})
	spec.Expect(filter("1234567", nil).(string)).ToEqual("123456")
}

func TestTruncateAStringFromAFloat(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{floatValue(3.3)})
	spec.Expect(filter("1234567", nil).(string)).ToEqual("123")
}

func TestTruncateAStringFromAString(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{stringValue("4")})
	spec.Expect(filter("1234567", nil).(string)).ToEqual("1234")
}

func TestTruncateAStringFromAnInvalidString(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{stringValue("abc")})
	spec.Expect(filter("1234567", nil).(string)).ToEqual("1234567")
}

func TestTruncateAnInvalidValue(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateFactory([]core.Value{intValue(4)})
	spec.Expect(filter(555, nil).(int)).ToEqual(555)
}
