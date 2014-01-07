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
	start := p.Position
	name := p.ReadName()
	if len(name) == 0 {
		return nil, p.Error("Invalid variable name in for tag", start)
	}
	if p.SkipSpaces() != 'i' || p.Left() < 3 || p.Data[p.Position+1] != 'n' || !core.IsTokenEnd(p.Data[p.Position+2]) {
		return nil, p.Error("Expecting keyword 'in' after variable name in for tag", start)
	}
	p.ForwardBy(2)

	value, err := p.ReadValue()
	if err != nil {
		return nil, err
	}

	p.SkipPastTag()
	return &For{
		NewCommon(),
		name,
		name + "[0]",
		name + "[1]",
		value,
	}, nil
}

func EndForFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endFor, nil
}

type For struct {
	*Common
	name      string
	keyName   string
	valueName string
	value     core.Value
}

func (f *For) AddSibling(tag core.Tag) error {
	return errors.New(fmt.Sprintf("%q does not belong inside of a for"))
}

func (f *For) Render(writer io.Writer, data map[string]interface{}) {
	resolved := f.value.Resolve(data)
	if s, ok := resolved.(string); ok {
		f.iterateString(writer, data, s)
	} else {
		value := reflect.ValueOf(resolved)
		kind := value.Kind()
		if kind == reflect.Array || kind == reflect.Slice {
			f.iterateArray(writer, data, value)
		} else if kind == reflect.Map {
			f.iterateMap(writer, data, value)
		}
	}
}

func (f *For) Name() string {
	return "for"
}

func (f *For) Type() core.TagType {
	return core.LoopTag
}

func (f *For) iterateString(writer io.Writer, data map[string]interface{}, s string) {
	length := len(s)
	if length == 0 {
		return
	}

	state := f.loopSetup(length, data, false)
	defer f.loopTeardown(state, data, false)

	for i := 0; i < length; i++ {
		f.loopIteration(state, i, length, data)
		data[f.name] = string(s[i])
		f.Common.Render(writer, data)
	}
}

func (f *For) iterateArray(writer io.Writer, data map[string]interface{}, value reflect.Value) {
	length := value.Len()
	if length == 0 {
		return
	}

	state := f.loopSetup(length, data, false)
	defer f.loopTeardown(state, data, false)

	for i := 0; i < length; i++ {
		f.loopIteration(state, i, length, data)
		data[f.name] = value.Index(i).Interface()
		f.Common.Render(writer, data)
	}
}

func (f *For) iterateMap(writer io.Writer, data map[string]interface{}, value reflect.Value) {
	length := value.Len()
	if length == 0 {
		return
	}

	state := f.loopSetup(length, data, true)
	defer f.loopTeardown(state, data, true)

	keys := value.MapKeys()
	for i := 0; i < length; i++ {
		f.loopIteration(state, i, length, data)
		key := keys[i]
		data[f.keyName] = key.Interface()
		data[f.valueName] = value.MapIndex(key).Interface()
		f.Common.Render(writer, data)
	}
}

func (f *For) loopSetup(length int, data map[string]interface{}, isMap bool) *State {
	state := &State{
		originalForLoop: data["forloop"],
		Length:          length,
	}
	data["forloop"] = state
	if isMap {
		state.key = data[f.keyName]
		state.value = data[f.valueName]
	} else {
		state.value = data[f.name]
	}
	return state
}

func (f *For) loopIteration(state *State, i, length int, data map[string]interface{}) {
	state.Index = i + 1
	state.Index0 = i
	state.RIndex = length - i
	state.RIndex0 = length - i - 1
	state.First = i == 0
	state.Last = i == length-1
}

func (f *For) loopTeardown(state *State, data map[string]interface{}, isMap bool) {
	data["forloop"] = state.originalForLoop
	if isMap {
		data[f.keyName] = state.key
		data[f.valueName] = state.value
	} else {
		data[f.name] = state.value
	}
}

type State struct {
	key             interface{}
	value           interface{}
	originalForLoop interface{}
	Length          int
	Index           int
	Index0          int
	RIndex          int
	RIndex0         int
	First           bool
	Last            bool
}
