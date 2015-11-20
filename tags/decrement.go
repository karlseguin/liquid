package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

// Creates an decrement tag
func DecrementFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	name := p.ReadName()
	if name == "" {
		name = globalIncrementlID
	}

	p.SkipPastTag()

	return &Decrement{
		name: name,
	}, nil
}

type Decrement struct {
	name string
}

func (d *Decrement) AddCode(code core.Code) {
	panic("Addcode should not have been called on a Decrement")
}

func (d *Decrement) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Decrement")
}

func (d *Decrement) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {

	// get count map for current data context
	countMap, ok := data[globalIncrementlID].(map[string]int)
	if !ok {
		countMap = map[string]int{}
	}
	count, _ := countMap[d.name]

	// decrement
	count = count - 1

	// get decrement string value
	valueString := core.ToString(count)

	// write
	writer.Write([]byte(valueString))

	countMap[d.name] = count
	data[globalIncrementlID] = countMap

	return core.Normal
}

func (d *Decrement) Name() string {
	return "decrement"
}

func (d *Decrement) Type() core.TagType {
	return core.StandaloneTag
}
