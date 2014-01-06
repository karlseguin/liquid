package core

import (
	"io"
)

// interface for something that can render itself
type Code interface {
	Render(writer io.Writer, data map[string]interface{})
}
