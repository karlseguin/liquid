// built in filters
package filters

// An interface function
type Filter func(input interface{}, data map[string]interface{}) interface{}

// a noop filter which returns the input as-is
// mostly used internally when the parameters don't make sense
func Noop(input interface{}, data map[string]interface{}) interface{} {
	return input
}
