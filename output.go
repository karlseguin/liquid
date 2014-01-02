package liquid

import (
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/filters"
)

func buildOutputTag(parser *core.Parser) (core.Token, error) {
	value, isStatic, err := parser.ReadValue()
	if err != nil {
		return nil, err
	}

	parser.SkipSpaces()
	var filters []filters.Filter
	if parser.Current() == '|' {
		parser.Forward()
		var err error
		if filters, err = buildFilters(parser); err != nil {
			return nil, err
		}
		filters = trimFilters(filters)
	}

	if isStatic {
		v := value.([]byte)
		if len(v) == 0 {
			return nil, nil
		}
		return &StaticOutput{Value: v, Filters: filters}, nil
	}
	return &DynamicOutput{Fields: value.([]string), Filters: filters}, nil
}




// 	//strip out leading and trailing {{ }}
// 	data := all[2 : len(all)-2]
// 	if start := core.SkipSpaces(data); start == -1 {
// 		return nil, nil
// 	} else {
// 		data = data[start:]
// 	}
// 	if data[0] == '\'' {
// 		static, position, err := createStaticOutput(data[1:], all)
// 		if err != nil {
// 			return nil, err
// 		}
// 		filters, err := extractFilters(data[position+2:], all)
// 		if err != nil {
// 			return nil, err
// 		}
// 		static.Filters = filters
// 		return static, nil
// 	}
// 	dynamic, position := createDynamicOutput(data, all)
// 	filters, err := extractFilters(data[position:], all)
// 	if err != nil {
// 		return nil, err
// 	}
// 	dynamic.Filters = filters
// 	return dynamic, nil
// }

// func extractFilters(data, all []byte) ([]filters.Filter, error) {
// 	filters := make([]filters.Filter, 0)
// 	for i, l := 0, len(data); i < l; i++ {
// 		b := data[i]
// 		if b == ' ' {
// 			continue
// 		}
// 		if b == '|' {
// 			filter, end, err := extracFilter(data[i+1:], all)
// 			if err != nil {
// 				return nil, err
// 			}
// 			filters = append(filters, filter)
// 			i += end
// 		} else {
// 			return nil, errors.New(fmt.Sprintf("invalid tag %q", all))
// 		}
// 	}
// 	return trimFilters(filters), nil
// }

// func extracFilter(data, all []byte) (filters.Filter, int, error) {
// 	start := core.SkipSpaces(data)
// 	if start == -1 {
// 		return nil, 0, errors.New(fmt.Sprintf("Empty filter in %q", all))
// 	}
// 	i := start
// 	l := len(data)
// 	for ; i < l; i++ {
// 		b := data[i]
// 		if b == ' ' || b == ':' {
// 			break
// 		}
// 	}
// 	name := strings.ToLower(string(data[start:i]))
// 	filterFactory, exists := Filters[name]
// 	if exists == false {
// 		return nil, 0, errors.New(fmt.Sprintf("Unknown filter %q in %q", name, all))
// 	}

// 	var parameters []string
// 	for ; i < l; i++ {
// 		b := data[i]
// 		if b == '|' {
// 			break
// 		}
// 		if b == ':' {
// 			p, position, err := extractParameters(data[i+1:], all)
// 			if err != nil {
// 				return nil, 0, err
// 			}
// 			i += position
// 			parameters = p
// 			break
// 		}
// 	}

// 	return filterFactory(parameters), i, nil
// }

// func extractParameters(data, all []byte) ([]string, int, error) {
// 	i := 0
// 	l := len(data)
// 	start := 0
// 	escaped := 0
// 	parameters := make([]string, 0, 1)
// 	for ; i < l; i++ {
// 		b := data[i]
// 		if b == ' ' || b == ',' {
// 			continue
// 		}
// 		if b == '|' {
// 			break
// 		}
// 		quoted := false
// 		if b == '\'' {
// 			quoted = true
// 			i++
// 		}
// 		start = i

// 		for j := start; j < l; j++ {
// 			b := data[j]
// 			if b == ',' || b == '|' || b == ' ' && quoted == false {
// 				parameters = append(parameters, string(data[start:j]))
// 				start = 0
// 				i = j
// 				break
// 			} else if quoted == true && b == '\'' {
// 				if data[j-1] == '\\' {
// 					escaped++
// 				} else {
// 					if escaped > 0 {
// 						parameters = append(parameters, string(unescape(data[start:j], escaped)))
// 					} else {
// 						parameters = append(parameters, string(data[start:j]))
// 					}
// 					start = 0
// 					escaped = 0
// 					i = j + 1
// 					break
// 				}
// 			}
// 		}
// 		if quoted == true && (data[i-1] != '\'' || start > 0) {
// 			return nil, 0, errors.New(fmt.Sprintf("Missing closing quote for parameter in %q", all))
// 		}
// 		if i == l-1 && start > 0 {
// 			if escaped > 0 {
// 				parameters = append(parameters, string(unescape(data[start:l], escaped)))
// 			} else {
// 				parameters = append(parameters, string(data[start:l]))
// 			}
// 			i = l
// 		}
// 	}
// 	return parameters, i, nil
// }

func trimFilters(values []filters.Filter) []filters.Filter {
	if len(values) == cap(values) {
		return values
	}
	trimmed := make([]filters.Filter, len(values))
	copy(trimmed, values)
	return trimmed
}
