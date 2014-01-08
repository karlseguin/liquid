package tags

import (
	"github.com/karlseguin/liquid/core"
	"io"
)

type Common struct {
	Code []core.Code
}

func NewCommon() *Common {
	return &Common{
		Code: make([]core.Code, 0, 5),
	}
}

func (c *Common) AddCode(code core.Code) {
	c.Code = append(c.Code, code)
}

func (c *Common) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	for _, code := range c.Code {
		if state := code.Execute(writer, data); state != core.Normal {
			return state
		}
	}
	return core.Normal
}
