package filters

import (
	"bytes"
	"github.com/karlseguin/liquid/core"
)

// Creates a capitalize filter
func CapitalizeFactory(parameters []core.Value) core.Filter {
	return Capitalize
}

// Capitalizes words in the input sentence
func Capitalize(input interface{}, data map[string]interface{}) interface{} {
	switch typed := input.(type) {
	case string:
		return capitalize([]byte(typed))
	case []byte:
		return capitalize(typed)
	default:
		return input
	}
}

// not crazy about this
func capitalize(sentence []byte) []byte {
	l := len(sentence)
	l1 := l - 1
	for i := 0; i < l1; i++ {
		if sentence[i] == ' ' && sentence[i+1] != ' ' {
			sentence[i+1] = bytes.ToUpper(sentence[i+1 : i+2])[0]
		}
	}
	if sentence[0] != ' ' {
		sentence[0] = bytes.ToUpper(sentence[0:1])[0]
	}
	return sentence
}
