// built in filters
package filters

type Filter func(input interface{}) interface{}

// a noop filter which returns the input as-is
// mostly used internally when the parameters don't make sense
func Noop(input interface{}) interface{} {
	return input
}
