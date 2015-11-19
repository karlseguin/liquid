package tags

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/core"
	"io"
	"reflect"
)

var (
	endTableRow = &End{"tablerow"}
)

type TableRowLoopState struct {
	data            map[string]interface{}
	writer          io.Writer
	key             interface{}
	value           interface{}
	originalForLoop interface{}
	offset          int
	isMap           bool
	Columns         int
	CurrentColumns  int
	CurrentRows     int
	Length          int
	Index           int
	Index0          int
	RIndex          int
	RIndex0         int
	First           bool
	Last            bool
}

func TableRowFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
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

	f := &TableRow{
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
		} else if name == "cols" {
			if p.SkipUntil(':') != ':' {
				return nil, p.Error("Expecting ':' after cols in for tag")
			}
			p.Forward()
			cols, err := p.ReadValue()
			if err != nil {
				return nil, err
			}
			f.cols = cols
		} else {
			return nil, p.Error(fmt.Sprint("%q is an inknown modifier in for tag", name))
		}
	}
	p.SkipPastTag()
	return f, nil
}

func EndTableRowFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endTableRow, nil
}

type TableRow struct {
	*Common
	name      string
	keyName   string
	valueName string
	cols      core.Value
	limit     core.Value
	offset    core.Value
	value     core.Value
}

func (f *TableRow) AddSibling(tag core.Tag) error {
	return errors.New(fmt.Sprintf("%q does not belong inside of a for", tag.Name()))
}

func (f *TableRow) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	value := reflect.ValueOf(f.value.Resolve(data))
	kind := value.Kind()

	if kind == reflect.Array || kind == reflect.Slice || kind == reflect.String || kind == reflect.Map {
		length := value.Len()
		if length == 0 {
			return core.Normal
		}

		state := &TableRowLoopState{
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
		if f.cols != nil {
			if cols, ok := core.ToInt(f.cols.Resolve(data)); ok {
				if cols < 1 {
					cols = 0
				}
				state.Columns = cols
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

func (f *TableRow) Name() string {
	return "tablerow"
}

func (f *TableRow) Type() core.TagType {
	return core.LoopTag
}

func (f *TableRow) iterateArray(state *TableRowLoopState, value reflect.Value, isString bool) {

	length := state.Length
	for i := 0; i < length; i++ {
		if state := f.iterateArrayIndex(i, state, value, isString); state == core.Break {
			return
		}
	}
}

func (f *TableRow) iterateArrayIndex(i int, state *TableRowLoopState, value reflect.Value, isString bool) core.ExecuteState {
	f.loopIteration(state, i)
	offsetI := i + state.offset
	item := value.Index(offsetI).Interface()
	if isString {
		item = string(item.(uint8))
	}
	state.data[f.name] = item

	return f.writeItem(state)
}

func (f *TableRow) iterateMap(state *TableRowLoopState, value reflect.Value) {
	f.writeRowStart(state)

	keys := value.MapKeys()
	length := state.Length
	for i := 0; i < length; i++ {
		if state := f.iterateMapIndex(i, state, keys, value); state == core.Break {
			return
		}
	}
	f.writeRowEnd(state)
}

func (f *TableRow) iterateMapIndex(i int, state *TableRowLoopState, keys []reflect.Value, value reflect.Value) core.ExecuteState {
	f.loopIteration(state, i)
	offsetI := i + state.offset
	key := keys[offsetI]
	state.data[f.keyName] = key.Interface()
	state.data[f.valueName] = value.MapIndex(key).Interface()
	return f.writeItem(state)
}

func (f *TableRow) loopIteration(state *TableRowLoopState, i int) {
	l1 := state.Length - 1
	state.Index = i + 1
	state.Index0 = i
	state.RIndex = state.Length - i
	state.RIndex0 = l1 - i
	state.First = i == l1
	state.Last = i == l1
}

func (f *TableRow) loopTeardown(state *TableRowLoopState) {
	data := state.data

	data["forloop"] = state.originalForLoop
	if state.isMap {
		data[f.keyName] = state.key
		data[f.valueName] = state.value
	} else {
		data[f.name] = state.value
	}
}

func (f *TableRow) writeRowStart(state *TableRowLoopState) {
	state.CurrentRows = state.CurrentRows + 1
	state.writer.Write([]byte(fmt.Sprintf("\n<tr class=\"row%v\">\n", state.CurrentRows)))
}
func (f *TableRow) writeRowEnd(state *TableRowLoopState) {
	state.writer.Write([]byte("</tr>\n"))
}
func (f *TableRow) writeItem(state *TableRowLoopState) core.ExecuteState {

	if state.Columns > 0 && state.CurrentColumns == 0 {
		f.writeRowStart(state)
	}
	state.CurrentColumns = state.CurrentColumns + 1

	state.writer.Write([]byte(fmt.Sprintf("<td class=\"col%v\">\n", state.CurrentColumns)))
	res := f.Common.Execute(state.writer, state.data)
	state.writer.Write([]byte(fmt.Sprintf("\n</td>\n")))

	if state.Columns > 0 && state.CurrentColumns >= state.Columns {
		state.CurrentColumns = 0
		f.writeRowEnd(state)
	}
	return res
}