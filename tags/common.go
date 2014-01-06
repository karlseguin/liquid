package tags

import (
	"bytes"
	"github.com/karlseguin/liquid/core"
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

func (c *Common) Render(data map[string]interface{}) []byte {
	buffer := new(bytes.Buffer)
	for _, code := range c.Code {
		buffer.Write(code.Render(data))
	}
	return buffer.Bytes()
}
