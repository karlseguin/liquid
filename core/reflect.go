package core

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	typeCache = make(map[reflect.Value]map[string]int)
	typeLock  sync.RWMutex
)

// Resolves the value of field within data
func Resolve(data interface{}, field string) interface{} {
	switch typed := data.(type) {
	case map[string]string:
		return typed[field]
	case map[string]interface{}:
		return typed[field]
	case map[string]int:
		return typed[field]
	case map[string]bool:
		return typed[field]
	case map[string]float64:
		return typed[field]
	case map[string]byte:
		return typed[field]
	case map[string][]byte:
		return typed[field]
	}
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(field))
		if value.IsValid() {
			return value.Interface()
		}
		return nil
	case reflect.Struct:
		return resolveStruct(v, field)
	case reflect.Ptr:
		v = reflect.Indirect(v)
		if v.Kind() != reflect.Struct {
			return nil
		}
		return resolveStruct(v, field)
	default:
		return nil
	}
}

// This is necessary in the case where the value references a struct directly:
//  template := "{{ user }}""
//  data := map[string]interface{}{"user": &User{"Leto"}},
// Without this step, the above would result in a value which points to the User
// we need to resolve this a step further and get the value of "User" (which
// will either me the output of its String() method, or %v)
//
// Of course, we only want this final resolution once we need the value. If we
// call this too early, say in Resolve above, we won't be able to build nested
// paths
func ResolveFinal(value interface{}) interface{} {
	if _, ok := value.(time.Time); ok {
		return value
	}
	kind := reflect.ValueOf(value).Kind()
	if kind == reflect.Ptr || kind == reflect.Struct {
		return resolvePtrOrStruct(value)
	}
	return value
}

func resolvePtrOrStruct(value interface{}) interface{} {
	if s, ok := value.(fmt.Stringer); ok {
		return s.String()
	}
	return ToBytes(value)
}

func resolveStruct(value reflect.Value, field string) interface{} {
	typeLock.RLock()
	typeData, exists := typeCache[value]
	typeLock.RUnlock()

	if exists == false {
		typeData = buildTypeData(value)
	}
	if index, exists := typeData[field]; exists {
		return value.Field(index).Interface()
	}
	return nil
}

func buildTypeData(value reflect.Value) map[string]int {
	t := value.Type()
	fieldCount := t.NumField()
	typeData := make(map[string]int, fieldCount)
	for i := 0; i < fieldCount; i++ {
		typeData[strings.ToLower(t.Field(i).Name)] = i
	}

	typeLock.Lock()
	defer typeLock.Unlock()
	if typeData, exists := typeCache[value]; exists {
		return typeData
	}
	typeCache[value] = typeData
	return typeData
}

//gets the length of string, map or array
func ToLength(input interface{}) (int, bool) {
	if s, ok := input.(string); ok {
		return len(s), true
	}

	value := reflect.ValueOf(input)
	kind := value.Kind()
	if kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map {
		return value.Len(), true
	}
	return 0, false
}
