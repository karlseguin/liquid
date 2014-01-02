package filters

func RemoveFactory(parameters []string) Filter {
	if len(parameters) != 1 {
		return Noop
	}
	return (&ReplaceFilter{parameters[0], "", -1}).Replace
}
