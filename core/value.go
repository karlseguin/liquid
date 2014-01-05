package core

// Represents a value, which can either be static or dynamic
type Value interface {
	Resolve(data map[string]interface{}) interface{}
	ResolveWithNil(data map[string]interface{}) interface{}
	Underlying() interface{}
}
