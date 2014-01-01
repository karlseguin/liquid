package liquid

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/filters"
	"github.com/karlseguin/liquid/helpers"
	"strings"
)

func outputExtractor(all []byte) (Token, error) {
	//strip out leading and trailing {{ }}
	data := all[2 : len(all)-2]
	if start := helpers.SkipSpaces(data); start == -1 {
		return nil, nil
	} else {
		data = data[start:]
	}
	if data[0] == '\'' {
		static, position, err := createStaticOutput(data[1:], all)
		if err != nil {
			return nil, err
		}
		filters, err := extractFilters(data[position+2:], all)
		if err != nil {
			return nil, err
		}
		static.Filters = filters
		return static, nil
	}
	dynamic, position := createDynamicOutput(data, all)
	filters, err := extractFilters(data[position:], all)
	if err != nil {
		return nil, err
	}
	dynamic.Filters = filters
	return dynamic, nil
}

func extractFilters(data, all []byte) ([]filters.Filter, error) {
	filters := make([]filters.Filter, 0, 0)
	for i, l := 0, len(data); i < l; i++ {
		b := data[i]
		if b == ' ' {
			continue
		}
		if b == '|' {
			filter, end, err := extracFilter(data[i+1:], all)
			if err != nil {
				return nil, err
			}
			filters = append(filters, filter)
			i += end
		} else {
			return nil, errors.New(fmt.Sprintf("invalid tag %q", all))
		}
	}
	return filters, nil
}

func extracFilter(data, all []byte) (filters.Filter, int, error) {
	start := helpers.SkipSpaces(data)
	if start == -1 {
		return nil, 0, errors.New(fmt.Sprintf("Empty filter in %q", all))
	}
	i := start
	l := len(data)
	for ; i < l; i++ {
		b := data[i]
		if b == ' ' || b == ':' {
			break
		}
	}
	name := strings.ToLower(string(data[start:i]))
	filterFactory, exists := Filters[name]
	if exists == false {
		return nil, 0, errors.New(fmt.Sprintf("Unknown filter %q in %q", name, all))
	}

	var parameters []string
	for ; i < l; i++ {
		b := data[i]
		if b == '|' {
			break
		}
		if b == ':' {
			p, position, err := extractParameters(data[i+1:], all)
			if err != nil {
				return nil, 0, err
			}
			i += position
			parameters = p
			break
		}
	}

	return filterFactory(parameters), i, nil
}

func extractParameters(data, all []byte) ([]string, int, error) {
	i := 0
	l := len(data)
	start := 0
	parameters := make([]string, 0, 1)

	for ; i < l; i++ {
		b := data[i]
		if b == ' ' || b == ',' {
			continue
		}
		if b == '|' {
			break
		}
		quoted := false
		if b == '\'' {
			quoted = true
			i++
		}
		start = i
		for j := start; j < l; j++ {
			b := data[j]

			if b == ',' || b == '|' || b == ' ' && quoted == false {
				parameters = append(parameters, string(data[start:j]))
				start = 0
				i = j
				break
			} else if b == '\'' && data[j-1] != '\\' && quoted == true {
				//todo unescape
				parameters = append(parameters, string(data[start:j]))
				start = 0
				i = j + 1
				break
			}
		}
		if quoted == true && (data[i-1] != '\'' || start > 0) {
			return nil, 0, errors.New(fmt.Sprintf("Missing closing quote for parameter in %q", all))
		}
		if i == l-1 && start > 0 {
			parameters = append(parameters, string(data[start:l]))
			i = l
		}
	}
	return parameters, i, nil
}
