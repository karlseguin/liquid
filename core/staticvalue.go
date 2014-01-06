package core

type StaticStringValue struct {
	Value string
}

func (v *StaticStringValue) Resolve(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticStringValue) ResolveWithNil(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticStringValue) Underlying() interface{} {
	return v.Value
}

type StaticIntValue struct {
	Value int
}

func (v *StaticIntValue) Resolve(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticIntValue) ResolveWithNil(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticIntValue) Underlying() interface{} {
	return v.Value
}

type StaticFloatValue struct {
	Value float64
}

func (v *StaticFloatValue) Resolve(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticFloatValue) ResolveWithNil(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticFloatValue) Underlying() interface{} {
	return v.Value
}

type StaticBoolValue struct {
	Value bool
}

func (v *StaticBoolValue) Resolve(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticBoolValue) ResolveWithNil(data map[string]interface{}) interface{} {
	return v.Value
}

func (v *StaticBoolValue) Underlying() interface{} {
	return v.Value
}

type StaticEmptyValue struct{}

func (v *StaticEmptyValue) Resolve(data map[string]interface{}) interface{} {
	return "liquid:empty"
}

func (v *StaticEmptyValue) ResolveWithNil(data map[string]interface{}) interface{} {
	return "liquid:empty"
}

func (v *StaticEmptyValue) Underlying() interface{} {
	return "liquid:empty"
}
