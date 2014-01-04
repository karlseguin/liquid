package core

// An interface function
type Filter func(input interface{}, data map[string]interface{}) interface{}

// A filter factory creates a filter based on the supplied parameters
type FilterFactory func(parameters []Value) Filter

// A map of filter names to filter factories
var FilterLookup = make(map[string]FilterFactory)

// Register's a filter for the given name (not thread-safe)
func RegisterFilter(name string, factory FilterFactory) {
	FilterLookup[name] = factory
}
