package filters

import (
	"github.com/karlseguin/liquid/core"
	"reflect"
	"sort"
)

// Creates a sort filter
func SortFactory(parameters []core.Value) core.Filter {
	return Sort
}

// Sorts an array
func Sort(input interface{}, data map[string]interface{}) interface{} {
	switch typed := input.(type) {
	case []int:
		sort.Ints(typed)
		return typed
	case []string:
		sort.Strings(typed)
		return typed
	case []float64:
		sort.Float64s(typed)
		return typed
	case sort.Interface:
		sort.Sort(typed)
		return typed
	}

	value := reflect.ValueOf(input)
	kind := value.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return input
	}
	l := value.Len()
	if l == 0 {
		return input
	}

	sortable := make(SortableObjects, l)
	for i := 0; i < l; i++ {
		value := value.Index(i).Interface()
		sortable[i] = &SortableObject{value, core.ToString(value)}
	}
	sort.Sort(sortable)

	sorted := make([]interface{}, l)
	for i := 0; i < l; i++ {
		sorted[i] = sortable[i].Underlying
	}
	return sorted
}

type SortableObjects []*SortableObject

func (s SortableObjects) Len() int {
	return len(s)
}

func (s SortableObjects) Less(i, j int) bool {
	return s[i].AsString < s[j].AsString
}

func (s SortableObjects) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SortableObject struct {
	Underlying interface{}
	AsString   string
}
