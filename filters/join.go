package filters

// Creates a join filter
func JoinFactory(parameters []string) Filter {
	j := &JoinFilter{
		glue: parameters[0],
	}
	return j.Join
}

type JoinFilter struct {
	glue string
}

// join elements of the array with certain character between them
func (f *JoinFilter) Join(input interface{}) interface{} {
  return nil
}
