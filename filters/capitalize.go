package filters

// Creates a capitalize filter
func CapitalizeFactory(parameters []string) Filter {
	return Capitalize
}

// Capitalizes words in the input sentence
func Capitalize(input interface{}) interface{} {
	return nil
}
