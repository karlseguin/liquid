package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
	"math/rand"
	"time"
)

var globalIncrementlID string

func init() {
	rand.Seed(time.Now().UnixNano())
	globalIncrementlID = "increment_" + randStringRunes(10)
}

// Creates an increment tag
func IncrementFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	name := p.ReadName()
	if name == "" {
		name = globalIncrementlID
	}

	p.SkipPastTag()

	return &Increment{
		name: name,
	}, nil
}

type Increment struct {
	name string
}

func (i *Increment) AddCode(code core.Code) {
	panic("Addcode should not have been called on a Increment")
}

func (i *Increment) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Increment")
}

func (i *Increment) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {

	// get count map for current data context
	countMap, ok := data[globalIncrementlID].(map[string]int)
	if !ok {
		countMap = map[string]int{}
	}
	count, _ := countMap[i.name]

	// get increment string value
	valueString := core.ToString(count)

	// write
	writer.Write([]byte(valueString))

	// increment
	countMap[i.name] = count + 1

	data[globalIncrementlID] = countMap

	return core.Normal
}

func (i *Increment) Name() string {
	return "increment"
}

func (i *Increment) Type() core.TagType {
	return core.StandaloneTag
}
