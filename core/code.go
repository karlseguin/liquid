package core

import (
	"io"
)

// interface for something that can render itself
type Code interface {
	Execute(writer io.Writer, data map[string]interface{}) ExecuteState
}

type ExecuteState int

const (
	Normal ExecuteState = iota
	Break
	Continue
)
