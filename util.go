package liquid

// return the position of the first none space, or -1 if no white space exists
func skipSpaces(data []byte) int {
	for index, b := range data {
		if b != ' ' {
			return index
		}
	}
	return -1
}

// Since these templates are possibly long-lived, let's free up any space
// which was accumulated while we grew these arrays
func trimArrayOfStrings(values []string) []string {
	if len(values) == cap(values) {
		return values
	}
	trimmed := make([]string, len(values))
	for index, value := range values {
		trimmed[index] = value
	}
	return trimmed
}
