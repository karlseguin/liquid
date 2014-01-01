package filters

// Creates a first filter
func FirstFactory(parameters []string) Filter {
	return First
}

// get the first element of the passed in array
func First(input interface{}) string {
  return ""
}
