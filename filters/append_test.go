package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestAppendToAString(t *testing.T) {
	spec := gspec.New(t)
	filter := AppendFactory([]core.Value{stringValue("?!")})
	spec.Expect(filter("dbz", nil).(string)).ToEqual("dbz?!")
}

func TestAppendToBytes(t *testing.T) {
	spec := gspec.New(t)
	filter := AppendFactory([]core.Value{stringValue("boring")})
	spec.Expect(filter([]byte("so"), nil).(string)).ToEqual("soboring")
}

func TestAppendADynamicValue(t *testing.T) {
	spec := gspec.New(t)
	filter := AppendFactory([]core.Value{dynamicValue("local.currency")})
	data := map[string]interface{}{
		"local": map[string]string{
			"currency": "$",
		},
	}
	spec.Expect(filter("100", data).(string)).ToEqual("100$")
}
