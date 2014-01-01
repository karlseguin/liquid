package liquid

import (
	"fmt"
	"errors"
)

func tagExtractor(all []byte) (Token, error) {
	//strip out leading and trailing {{ }}
	data := all[2:len(all)-2]
	if start := skipSpaces(data); start == -1 {
		return nil, nil
	} else {
		data = data[start:]
	}
	var token Token
	var err error
	if data[0] == '\'' {
		token, err = createStaticTag(data[1:], all)
	} else {
		token, err = createDynamicTag(data, all)
	}
	if err != nil {
		return nil, err
	}
	return token, nil
}

func createStaticTag(data, all []byte) (*StaticTag, error) {
	escaped := 0
	escaping := false
	for index, b := range data {
		if b == '\'' {
			if escaping {
				escaped++
				escaping = false
			} else {
				var value []byte
				if escaped > 0 {
					return &StaticTag{Value: unescape(data[0:index], escaped)}, nil
				}
				value = make([]byte, index)
				copy(value, data[:index])
				return &StaticTag{Value: value}, nil
			}
		} else if b == '\\' && escaping == false {
			escaping = true
		} else {
			escaping = false
		}
	}
	return nil, errors.New(fmt.Sprintf("Tag had an unclosed single quote in %q", all))
}

func unescape(data []byte, escaped int) []byte {
	value := make([]byte, len(data) - escaped)
	i := 0
	found := 0
	position := 0
	for l := len(data)-1; i < l; i++ {
		b := data[i]
		if b == '\\' && data[i+1] == '\'' {
			value[position] = '\''
			found++
			i++
			if found == escaped { break }
		} else {
			value[position] = b
		}
		position++
	}
	copy(value[position:], data[i:])
	return value
}

func createDynamicTag(data, all []byte) (*DynamicTag, error) {
	start := 0
	values := make([][]byte, 0, 1)
	for index, b := range data {
		if b == ' ' {
			values = append(values, data[start:index])
			break
		}
		if b == '.' {
			values = append(values, data[start:index])
			start = index + 1
		}
	}
	return &DynamicTag{
		Values: trimArrayOfBytes(values),
	}, nil
}
