package filters

// Creates a last filter
func LastFactory(parameters []string) Filter {
	return Last
}

// get the last element of the passed in array
func Last(input interface{}) string {
  return ""
}
