package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var globalCyclelID string

func init() {
	rand.Seed(time.Now().UnixNano())
	globalCyclelID = "cycle_" + randStringRunes(10)
}

// Creates an cycle tag
func CycleFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	name, values, _ := p.ReadCycleValues()
	if len(values) == 0 {
		return nil, p.Error("Invalid cycle values")
	}
	p.SkipPastTag()

	if name == "" {
		name = globalCyclelID
	}

	return &Cycle{
		name:   name,
		values: values,
	}, nil
}

type Cycle struct {
	name   string
	values []core.Value
}

func (c *Cycle) AddCode(code core.Code) {
	panic("Addcode should not have been called on a Cycle")
}

func (c *Cycle) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Cycle")
}

func (c *Cycle) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {

	// get count map for current data context
	countMap, ok := data[globalCyclelID].(map[string]int)
	if !ok {
		countMap = map[string]int{}
	}
	idx, _ := countMap[c.name]
	if idx >= len(c.values) {
		countMap[c.name] = 0
		idx = 0
	}

	// get cycle string value
	value := c.values[idx]
	valueString := core.ToString(value.Resolve(nil))

	// write
	writer.Write([]byte(valueString))

	// increment
	countMap[c.name] = idx + 1

	data[globalCyclelID] = countMap

	return core.Normal
}

func (c *Cycle) Name() string {
	return "cycle"
}

func (c *Cycle) Type() core.TagType {
	return core.StandaloneTag
}
