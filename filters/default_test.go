package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestDefaultWithBuiltinValue(t *testing.T) {
	spec := gspec.New(t)
	filter := DefaultFactory(nil)
	spec.Expect(string(filter(nil, nil).([]byte))).ToEqual("")
}

func TestDefaultWithValueOnString(t *testing.T) {
	spec := gspec.New(t)
	filter := DefaultFactory([]core.Value{stringValue("d")})
	spec.Expect(string(filter("", nil).([]byte))).ToEqual("d")
}

func TestDefaultWithValueOnArray(t *testing.T) {
	spec := gspec.New(t)
	filter := DefaultFactory([]core.Value{stringValue("n/a")})
	spec.Expect(string(filter([]int{}, nil).([]byte))).ToEqual("n/a")
}
