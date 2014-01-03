package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
)

func TestModuloAnIntToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := ModuloFactory([]core.Value{intValue(5)})
	spec.Expect(filter(43, nil).(int)).ToEqual(3)
}

func TestModuloAnFloatToAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := ModuloFactory([]core.Value{floatValue(5.2)})
	spec.Expect(filter(43, nil).(int)).ToEqual(3)
}

func TestModuloAnIntToAStringAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := ModuloFactory([]core.Value{intValue(7)})
	spec.Expect(filter("33", nil).(int)).ToEqual(5)
}

func TestModuloAnIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := ModuloFactory([]core.Value{intValue(7)})
	spec.Expect(filter([]byte("34"), nil).(int)).ToEqual(6)
}

func TestModuloAnDynamicIntToBytesAsAnInt(t *testing.T) {
	spec := gspec.New(t)
	filter := ModuloFactory([]core.Value{dynamicValue("fee")})
	spec.Expect(filter([]byte("34"), params("fee", 5)).(int)).ToEqual(4)
}
