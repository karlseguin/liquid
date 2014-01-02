// shared interfaces and utility functions
package core

type Tag interface {
	Token
	Container
	Name() string
	IsEnd() bool
}
