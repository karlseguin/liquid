package helpers

import (
	"strconv"
	"fmt"
)

// return the position of the first none space, or -1 if no white space exists
func SkipSpaces(data []byte) int {
	for index, b := range data {
		if b != ' ' {
			return index
		}
	}
	return -1
}

// Since these templates are possibly long-lived, let's free up any space
// which was accumulated while we grew these arrays
func TrimArrayOfStrings(values []string) []string {
	if len(values) == cap(values) {
		return values
	}
	trimmed := make([]string, len(values))
	for index, value := range values {
		trimmed[index] = value
	}
	return trimmed
}

// Convert arbitrary data to []byte
func DataToBytes(data interface{}) []byte {
	switch typed := data.(type) {
	case []byte:
		return typed
	case string:
		return []byte(typed)
	case bool:
		return []byte(strconv.FormatBool(typed))
	case float64:
		return []byte(strconv.FormatFloat(typed, 'g', -1, 64))
	case uint64:
		return []byte(strconv.FormatUint(typed, 10))
	case uint:
		return []byte(strconv.FormatUint(uint64(typed), 10))
	case int:
		return []byte(strconv.Itoa(typed))
	case fmt.Stringer:
		return []byte(typed.String())
	}
	return []byte(fmt.Sprintf("%v", data))
}
