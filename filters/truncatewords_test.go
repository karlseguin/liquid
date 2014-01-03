package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestTruncateWordsWhenTooLong(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateWordsFactory([]core.Value{intValue(2)})
	spec.Expect(filter("hello world how's it going", nil).(string)).ToEqual("hello world...")
}

func TestTruncateWordsWhenTooShort(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateWordsFactory([]core.Value{intValue(6)})
	spec.Expect(filter("hello world how's it going", nil).(string)).ToEqual("hello world how's it going")
}

func TestTruncateWordsWhenJustRight(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateWordsFactory([]core.Value{intValue(5)})
	spec.Expect(filter("hello world how's it going", nil).(string)).ToEqual("hello world how's it going")
}

func TestTruncateWordsWithCustomAppend(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateWordsFactory([]core.Value{intValue(3), stringValue(" (more...)")})
	spec.Expect(filter("hello world how's it going", nil).(string)).ToEqual("hello world how's (more...)")
}

func TestTruncateWordsWithShortWords(t *testing.T) {
	spec := gspec.New(t)
	filter := TruncateWordsFactory([]core.Value{dynamicValue("max")})
	spec.Expect(filter("I  think  a  feature good", params("max", 2)).(string)).ToEqual("I  think  a  feature...")
}
