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

func (c *Common) Render(writer io.Writer, data map[string]interface{}) {
	for _, code := range c.Code {
		code.Render(writer, data)
	}
}
