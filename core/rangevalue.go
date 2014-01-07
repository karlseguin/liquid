package core

type RangeValue struct {
	from Value
	to   Value
}

func (v *RangeValue) ResolveWithNil(data map[string]interface{}) interface{} {
	return v.Resolve(data)
}

func (v *RangeValue) Resolve(data map[string]interface{}) interface{} {
	from, to := 0, 0
	ok := false
	if from, ok = ToInt(v.from.Resolve(data)); ok == false {
		return nil
	}
	if to, ok = ToInt(v.to.Resolve(data)); ok == false {
		return nil
	}
	length := to - from + 1
	if length < 1 {
		return nil
	}
	n := make([]int, length)
	for i := 0; i < length; i++ {
		n[i] = i + 1
	}
	return n
}

func (v *RangeValue) Underlying() interface{} {
	return nil
}
