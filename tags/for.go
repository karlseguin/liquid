package tags

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/core"
	"io"
	"reflect"
)

var (
	endFor = &End{"for"}
)

func ForFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	name := p.ReadName()
	if len(name) == 0 {
		return nil, p.Error("Invalid variable name in for tag")
	}
	if p.SkipSpaces() != 'i' || p.Left() < 3 || p.Data[p.Position+1] != 'n' || !core.IsTokenEnd(p.Data[p.Position+2]) {
		return nil, p.Error("Expecting keyword 'in' after variable name in for tag")
	}
	p.ForwardBy(2)

	value, err := p.ReadValue()
	if err != nil {
		return nil, err
	}

	f := &For{
		Common:    NewCommon(),
		name:      name,
		keyName:   name + "[0]",
		valueName: name + "[1]",
		value:     value,
	}

	for {
		name := p.ReadName()
		if name == "" {
			break
		}
		if name == "limit" {
			if p.SkipUntil(':') != ':' {
				return nil, p.Error("Expecting ':' after limit in for tag")
			}
			p.Forward()
			limit, err := p.ReadValue()
			if err != nil {
				return nil, err
			}
			f.limit = limit
		} else if name == "offset" {
			if p.SkipUntil(':') != ':' {
				return nil, p.Error("Expecting ':' after offset in for tag")
			}
			p.Forward()
			offset, err := p.ReadValue()
			if err != nil {
				return nil, err
			}
			f.offset = offset
		} else if name == "reverse" {
			f.reverse = true
		} else {
			return nil, p.Error(fmt.Sprint("%q is an inknown modifier in for tag", name))
		}
	}
	p.SkipPastTag()
	return f, nil
}

func EndForFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endFor, nil
}

type For struct {
	*Common
	name      string
	keyName   string
	valueName string
	reverse   bool
	limit     core.Value
	offset    core.Value
	value     core.Value
}

func (f *For) AddSibling(tag core.Tag) error {
	return errors.New(fmt.Sprintf("%q does not belong inside of a for", tag.Name()))
}

func (f *For) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	value := reflect.ValueOf(f.value.Resolve(data))
	kind := value.Kind()

	if kind == reflect.Array || kind == reflect.Slice || kind == reflect.String || kind == reflect.Map {
		length := value.Len()
		if length == 0 {
			return core.Normal
		}

		state := &LoopState{
			data:            data,
			writer:          writer,
			originalForLoop: data["forloop"],
			Length:          length,
		}
		data["forloop"] = state
		if f.limit != nil {
			if limit, ok := core.ToInt(f.limit.Resolve(data)); ok && limit < length {
				state.Length = limit
			}
		}

		if f.offset != nil {
			if offset, ok := core.ToInt(f.offset.Resolve(data)); ok {
				state.offset = offset
				if n := state.Length + state.offset; n < state.Length {
					state.Length = n
				} else if length < state.offset+state.Length {
					state.Length = length - offset
				}
			}
		}

		defer f.loopTeardown(state)
		if kind == reflect.Map {
			state.isMap = true
			state.key = data[f.keyName]
			state.value = data[f.valueName]
			f.iterateMap(state, value)
		} else {
			state.value = data[f.name]
			f.iterateArray(state, value, kind == reflect.String)
		}
	}
	return core.Normal
}

func (f *For) Name() string {
	return "for"
}

func (f *For) Type() core.TagType {
	return core.LoopTag
}

func (f *For) iterateArray(state *LoopState, value reflect.Value, isString bool) {
	length := state.Length
	if f.reverse {
		for i := length - 1; i >= 0; i-- {
			if state := f.iterateArrayIndex(i, state, value, isString); state == core.Break {
				return
			}
		}
	} else {
		for i := 0; i < length; i++ {
			if state := f.iterateArrayIndex(i, state, value, isString); state == core.Break {
				return
			}
		}
	}
}

func (f *For) iterateArrayIndex(i int, state *LoopState, value reflect.Value, isString bool) core.ExecuteState {
	f.loopIteration(state, i)
	offsetI := i + state.offset
	item := value.Index(offsetI).Interface()
	if isString {
		item = string(item.(uint8))
	}
	state.data[f.name] = item
	return f.Common.Execute(state.writer, state.data)
}

func (f *For) iterateMap(state *LoopState, value reflect.Value) {
	keys := value.MapKeys()
	length := state.Length
	if f.reverse {
		for i := length - 1; i >= 0; i-- {
			if state := f.iterateMapIndex(i, state, keys, value); state == core.Break {
				return
			}
		}
	} else {
		for i := 0; i < length; i++ {
			if state := f.iterateMapIndex(i, state, keys, value); state == core.Break {
				return
			}
		}
	}
}

func (f *For) iterateMapIndex(i int, state *LoopState, keys []reflect.Value, value reflect.Value) core.ExecuteState {
	f.loopIteration(state, i)
	offsetI := i + state.offset
	key := keys[offsetI]
	state.data[f.keyName] = key.Interface()
	state.data[f.valueName] = value.MapIndex(key).Interface()
	return f.Common.Execute(state.writer, state.data)
}

func (f *For) loopIteration(state *LoopState, i int) {
	l1 := state.Length - 1
	if f.reverse {
		state.Index = state.Length - i
		state.Index0 = l1 - i
		state.RIndex = i + 1
		state.RIndex0 = i
		state.First = i == l1
		state.Last = i == 0
	} else {
		state.Index = i + 1
		state.Index0 = i
		state.RIndex = state.Length - i
		state.RIndex0 = l1 - i
		state.First = i == l1
		state.Last = i == l1
	}
}

func (f *For) loopTeardown(state *LoopState) {
	data := state.data
	data["forloop"] = state.originalForLoop
	if state.isMap {
		data[f.keyName] = state.key
		data[f.valueName] = state.value
	} else {
		data[f.name] = state.value
	}
}

type LoopState struct {
	data            map[string]interface{}
	writer          io.Writer
	key             interface{}
	value           interface{}
	originalForLoop interface{}
	offset          int
	isMap           bool
	Length          int
	Index           int
	Index0          int
	RIndex          int
	RIndex0         int
	First           bool
	Last            bool
}
